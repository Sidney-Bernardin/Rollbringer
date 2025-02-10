package database

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"maps"
	"slices"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
)

var database *sqlx.DB

type DatabaseRepository struct {
	DB *sqlx.DB
	TX sqlx.ExtContext
}

func NewDatabase(config *domain.Config, migrations fs.FS) (*DatabaseRepository, error) {
	repo := &DatabaseRepository{
		DB: database,
		TX: database,
	}

	if database != nil {
		return repo, nil
	}

	// Create connection to Postgres server.
	db, err := sqlx.Connect("pgx", config.PGUrl)
	if err != nil {
		return nil, domain.Wrap(err, "cannot connect to Postgres server", nil)
	}
	database = db
	repo.DB = db
	repo.TX = db

	// Create migration source.
	migrationSrc, err := iofs.New(migrations, "migrations")
	if err != nil {
		return nil, domain.Wrap(err, "cannot create migration source", nil)
	}

	// Create a migrate instance.
	migrator, err := migrate.NewWithSourceInstance("iofs", migrationSrc, config.PGUrl)
	if err != nil {
		return nil, domain.Wrap(err, "cannot create migrate instance", nil)
	}

	// Migrate database all the way up.
	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, domain.Wrap(err, "cannot migrate all the way up", nil)
	}

	return repo, nil
}

func (repo *DatabaseRepository) Close() error {
	err := repo.DB.Close()
	return domain.Wrap(err, "cannot close Postgres connection", nil)
}

func (repo *DatabaseRepository) Insert(ctx context.Context, record any, query string, args ...any) error {
	if err := sqlx.GetContext(ctx, repo.TX, record, query, args...); err != nil {
		if pgErr, ok := errors.Cause(err).(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return domain.ErrAlreadyExists
		}

		return domain.Wrap(err, "cannot insert record", nil)
	}

	return nil
}

func (repo *DatabaseRepository) Get(ctx context.Context, record any, query string, args ...any) error {
	if err := sqlx.GetContext(ctx, repo.TX, record, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrNotFound
		}

		return domain.Wrap(err, "cannot get record", nil)
	}

	return nil
}

func (repo *DatabaseRepository) Update(ctx context.Context, updates map[string]any, query string, args ...any) error {

	var sets string
	for n, k := range slices.Collect(maps.Keys(updates)) {
		sets += fmt.Sprintf("%s = $%d,", k, n+1+len(args))
	}
	sets = sets[:len(sets)-1]

	result, err := repo.TX.ExecContext(ctx, fmt.Sprintf(query, sets),
		append(args, slices.Collect(maps.Values(updates))...)...)

	if err != nil {
		return domain.Wrap(err, "cannot update record", nil)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return domain.Wrap(err, "cannot get affected rows", nil)
	}

	if affected <= 0 {
		return domain.ErrNotFound
	}

	return nil
}

package database

import (
	"context"
	"database/sql"
	"embed"
	"rollbringer/src/domain"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Database struct {
	DB *sqlx.DB
	TX sqlx.ExtContext
}

func NewDatabase(dbURL string, migrations *embed.FS) (*Database, error) {

	// Connect to Postgres.
	db, err := sqlx.Connect("pgx", dbURL)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to Postgres server")
	}

	// Create migration source.
	migrationSrc, err := iofs.New(migrations, "migrations")
	if err != nil {
		return nil, errors.Wrap(err, "cannot create migration source")
	}

	// Create migrate instance.
	migrator, err := migrate.NewWithSourceInstance("iofs", migrationSrc, dbURL)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create migrate instance")
	}

	// Migrate database.
	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, errors.Wrap(err, "cannot migrate")
	}

	return &Database{db, db}, nil
}

func (db *Database) Close() error {
	err := db.DB.Close()
	return errors.Wrap(err, "cannot close database")
}

type CRUDFunc func(ctx context.Context, view any, q string, args ...any) error

func (db *Database) CRUDInsert(ctx context.Context, view any, query string, args ...any) (err error) {
	if view == nil {
		_, err = db.TX.ExecContext(ctx, query, args...)
	} else {
		err = sqlx.GetContext(ctx, db.TX, view, query, args...)
	}

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return domain.ErrEntityConflict
			}
		}

		return errors.Wrap(err, "cannot insert row")
	}

	return nil
}

func (db *Database) CRUDGet(ctx context.Context, view any, query string, args ...any) error {
	if err := sqlx.GetContext(ctx, db.TX, view, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrEntityNotFound
		}

		return errors.Wrap(err, "cannot get row")
	}

	return nil
}

func (db *Database) CRUDGetMany(ctx context.Context, view any, query string, args ...any) error {
	if err := sqlx.SelectContext(ctx, db.TX, view, query, args...); err != nil {
		return errors.Wrap(err, "cannot get rows")
	}
	return nil
}

func (db *Database) CRUDUpdate(ctx context.Context, view any, query string, args ...any) error {
	result, err := db.TX.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "cannot update row")
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "cannot get affected rows")
	}

	if affected <= 0 {
		return domain.ErrNoEntitiesEffected
	}

	return nil
}

package sql

import (
	"context"
	"database/sql"
	"embed"
	"log/slog"

	"github.com/Sidney-Bernardin/Rollbringer/internal"
	"github.com/Sidney-Bernardin/Rollbringer/internal/domain"

	gomigrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const version = 20250530153240

//go:embed Migrations/*.sql
var migrations embed.FS

type SQL interface {
	GetUser(ctx context.Context, userID domain.UUID) (*domain.User, error)
}

type repository struct {
	config *internal.Config
	log    *slog.Logger

	db *sqlx.DB
	tx sqlx.ExtContext
}

func New(ctx context.Context, config *internal.Config, log *slog.Logger) (SQL, error) {

	db, err := sqlx.ConnectContext(ctx, "pgx", config.PostgresURL)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to Postgres")
	}

	if err := migrate(config.PostgresURL); err != nil {
		return nil, errors.Wrap(err, "cannot do migrate")
	}

	return &repository{config, log, db, db}, nil
}

func migrate(url string) error {

	migrationSrc, err := iofs.New(migrations, "Migrations")
	if err != nil {
		return errors.Wrap(err, "cannot create migration source")
	}

	m, err := gomigrate.NewWithSourceInstance("iofs", migrationSrc, url)
	if err != nil {
		return errors.Wrap(err, "cannot create migrate instance")
	}
	defer m.Close()

	if err := m.Migrate(version); err != nil && err != gomigrate.ErrNoChange {
		return errors.Wrap(err, "cannot migrate")
	}

	return nil
}

func (repo *repository) transaction(ctx context.Context, callback func(tx *repository) error) error {

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "cannot begin transaction")
	}

	defer func() {
		repo.log.Warn("Transaction rollback failed", "err", err.Error())
	}()

	txRepo := &repository{
		config: repo.config,
		log:    repo.log,
		tx:     tx,
	}

	if err := callback(txRepo); err != nil {
		return errors.Wrap(err, "transaction failed")
	}

	err = tx.Commit()
	return errors.Wrap(err, "cannot commit transaction")
}

func parseInsertErr[E any](err error) error {

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return &domain.EntityConflictError[E]{
			Column:  pgErr.ColumnName,
			Message: err.Error(),
		}
	}

	return err
}

func parseGetErr[E any](err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return &domain.NoEntitiesError[E]{}
	}
	return err
}

func parseUpdateErr[E any](result sql.Result, err error) error {
	if err != nil {
		return err
	}

	if updates, _ := result.RowsAffected(); updates <= 0 {
		return &domain.NoEntitiesEffectedError[E]{}
	}

	return nil
}

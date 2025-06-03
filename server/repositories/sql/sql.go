package sql

import (
	"context"
	"embed"
	"log/slog"

	"github.com/Sidney-Bernardin/Rollbringer/server"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/sql/queries"

	gomigrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

const version = 20250530153240

//go:embed Migrations/*.sql
var migrations embed.FS

var ErrNoRows = pgx.ErrNoRows

type SQL struct {
	*queries.Queries

	config *server.Config
	log    *slog.Logger

	pool *pgxpool.Pool
}

func New(ctx context.Context, config *server.Config, log *slog.Logger) (*SQL, error) {

	pool, err := pgxpool.New(ctx, config.PostgresUrl)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create create connection pool")
	}

	if err := migrate(config.PostgresUrl); err != nil {
		return nil, errors.Wrap(err, "cannot migrate")
	}

	return &SQL{
		Queries: queries.New(pool),
		config:  config,
		log:     log,
		pool:    pool,
	}, nil
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

func (sql *SQL) Transaction(ctx context.Context, callback func(tx *SQL) error) error {

	tx, err := sql.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return errors.Wrap(err, "cannot begin transaction")
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			sql.log.Warn("Transaction rollback failed", "err", err.Error())
		}
	}()

	txRepo := &SQL{
		Queries: sql.Queries.WithTx(tx),
		config:  sql.config,
		log:     sql.log,
	}

	if err := callback(txRepo); err != nil {
		return errors.Wrap(err, "transaction failed")
	}

	err = tx.Commit(ctx)
	return errors.Wrap(err, "cannot commit transaction")
}

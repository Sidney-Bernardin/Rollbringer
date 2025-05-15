package database

import (
	"context"
	"embed"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
)

//go:embed Migrations/*.sql
var migrations embed.FS

type Database struct {
	Pool *pgxpool.Pool
}

func NewDatabase(ctx context.Context, dbURL string) (*Database, error) {

	// Connect to the Postgres server.
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to Postgres server")
	}

	// Create migration source.
	migrationSrc, err := iofs.New(migrations, "Migrations")
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

	return &Database{pool}, nil
}

func (db *Database) Close() {
	db.Pool.Reset()
}

func (db *Database) Transaction(ctx context.Context, txFunc func(tx pgx.Tx) error) error {

	// Begin transaction.
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot begin transaction")
	}
	defer tx.Rollback(ctx)

	// Do transaction callback.
	if err := txFunc(tx); err != nil {
		return errors.Wrap(err, "transaction failed")
	}

	// Commit transaction.
	err = tx.Commit(ctx)
	return errors.Wrap(err, "cannot commit transaction")
}

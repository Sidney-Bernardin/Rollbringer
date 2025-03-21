package database

import (
	"embed"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"rollbringer/src"
)

var database *Database

type Database struct {
	DB *sqlx.DB
	TX sqlx.ExtContext
}

func NewDatabase(config *src.Config, migrations *embed.FS) (*Database, error) {
	if database != nil {
		return database, nil
	}

	// Connect to Postgres.
	db, err := sqlx.Connect("pgx", config.PostgresURL)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to Postgres server")
	}

	// Create migration source.
	migrationSrc, err := iofs.New(migrations, "migrations")
	if err != nil {
		return nil, errors.Wrap(err, "cannot create migration source")
	}

	// Create migrate instance.
	migrator, err := migrate.NewWithSourceInstance("iofs", migrationSrc, config.PostgresURL)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create migrate instance")
	}

	// Migrate database.
	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, errors.Wrap(err, "cannot migrate")
	}

	database = &Database{
		DB: db,
		TX: db,
	}

	return database, nil
}

func (db *Database) Close() error {
	err := db.DB.Close()
	return errors.Wrap(err, "cannot close database")
}

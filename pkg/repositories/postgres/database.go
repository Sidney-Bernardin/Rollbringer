package postgres

import (
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"rollbringer/pkg/domain"
)

var db *sqlx.DB

func NewDatabaseConnection(config *domain.Config, migrations fs.FS) (*sqlx.DB, error) {
	if db != nil {
		return db, nil
	}

	// Create connection to Postgres server.
	db, err := sqlx.Connect("pgx", config.PGUrl)
	if err != nil {
		return nil, domain.Wrap(err, "cannot connect to Postgres server", nil)
	}

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

	return db, nil
}

type DatabaseRepository struct {
	DB *sqlx.DB
	TX sqlx.ExtContext
}

func (repo *DatabaseRepository) Close() error {
	err := repo.DB.Close()
	return domain.Wrap(err, "cannot close Postgres connection", nil)
}

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

type DatabaseRepository struct {
	DB *sqlx.DB
	TX sqlx.ExtContext
}

func NewDatabaseRepository(config *domain.Config, migrations fs.FS) (ret *DatabaseRepository, err error) {

	// Create connection to Postgres server.
	db, err = sqlx.Connect("pgx", config.PGUrl)
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

	return &DatabaseRepository{db, db}, nil
}

func (repo *DatabaseRepository) Close() error {
	err := repo.DB.Close()
	return domain.Wrap(err, "cannot close Postgres connection", nil)
}

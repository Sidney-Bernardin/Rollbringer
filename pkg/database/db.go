package database

import (
	"database/sql"
	"os"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type Database struct {
	db *sql.DB
}

func New(addr string) (*Database, error) {

	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open database")
	}

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "ping failed")
	}

	database := &Database{db}
	if err := database.applySchema(); err != nil {
		return nil, errors.Wrap(err, "cannot apply schema")
	}

	return database, nil
}

func (database *Database) applySchema() error {

	b, err := os.ReadFile("schema.sql")
	if err != nil {
		return errors.Wrap(err, "cannot read schema.sql")
	}

	_, err = database.db.Exec(string(b))
	return errors.Wrap(err, "cannot execute schema")
}

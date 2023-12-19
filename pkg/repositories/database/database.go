package database

import (
	"database/sql"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var (
	ErrSessionNotFound = errors.New("session was not found")
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

	return &Database{db}, nil
}

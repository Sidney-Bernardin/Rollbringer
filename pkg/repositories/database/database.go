package database

import (
	"context"
	"slices"
	"strings"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type Database struct {
	tx sqlx.ExtContext
}

// New returns a new Database that connects to a Postgres server.
func New(addr string) (*Database, error) {

	// Create connection to the Postgres server.
	db, err := sqlx.Open("postgres", addr)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open postgres connection")
	}

	// Ping the Postgres server.
	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "ping failed")
	}

	return &Database{db}, nil
}

func (db *Database) Transaction(ctx context.Context, txFunc func(db *Database) error) error {

	tx, err := db.tx.(*sqlx.DB).Beginx()
	if err != nil {
		return errors.Wrap(err, "cannot begin transaction")
	}

	if err := txFunc(&Database{tx}); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return errors.Wrap(rbErr, "cannot rollback failed transaction")
		}

		return errors.Wrap(err, "transaction failed")
	}

	err = tx.Commit()
	return errors.Wrap(err, "cannot commit transaction")
}

func parseColumns(columns ...string) string {
	columns = slices.DeleteFunc(columns, func(c string) bool {
		return c == ""
	})
	return strings.Join(columns, ", ")
}

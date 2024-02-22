package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Database struct {
	conn *pgxpool.Pool
}

// New returns a new Database that connects to a Postgres server.
func New(addr string) (*Database, error) {

	// Connect to the Postgres server.
	conn, err := pgxpool.New(context.Background(), addr)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to postgres server")
	}

	// Ping the Postgres server.
	if err := conn.Ping(context.Background()); err != nil {
		return nil, errors.Wrap(err, "cannot ping postgres server")
	}

	return &Database{conn}, nil
}

package database

import (
	"context"
	"database/sql"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	ErrUnauthorized = errors.New("unauthorized")

	ErrUserNotFound = errors.New("user was not found")

	ErrMaxGames     = errors.New("max games reached")
	ErrGameNotFound = errors.New("game was not found")
)

type Database struct {
	db     *sql.DB
	logger *zerolog.Logger
}

// New returns a new Database that connects to the given Postgres server.
func New(addr string, loggerCtx context.Context) (*Database, error) {

	// Connect to the Postgres server.
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open database")
	}

	// Ping the Postgres server.
	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "ping failed")
	}

	return &Database{
		db:     db,
		logger: log.Ctx(loggerCtx),
	}, nil
}

// rollback rolls back the given transaction.
func (database *Database) rollback(tx *sql.Tx) {
	if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
		database.logger.Error().Stack().Err(err).Msg("Cannot rollback transaction")
	}
}

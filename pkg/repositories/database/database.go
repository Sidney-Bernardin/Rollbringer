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
	ErrUserNotFound    = errors.New("user was not found")
	ErrSessionNotFound = errors.New("session was not found")

	ErrMaxGames     = errors.New("max games reached")
	ErrGameNotFound = errors.New("game was not found")
)

type Database struct {
	db     *sql.DB
	logger *zerolog.Logger
}

func New(addr string, loggerCtx context.Context) (*Database, error) {

	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open database")
	}

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "ping failed")
	}

	return &Database{
		db:     db,
		logger: log.Ctx(loggerCtx),
	}, nil
}

func (database *Database) rollback(tx *sql.Tx) {
	if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
		database.logger.Error().Stack().Err(err).Msg("Cannot rollback transaction")
	}
}

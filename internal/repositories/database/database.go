package database

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"rollbringer/internal/config"
)

type Database struct {
	Config *config.Config
	Logger *slog.Logger

	TX sqlx.ExtContext
}

func New(cfg *config.Config, logger *slog.Logger) (*Database, error) {
	db, err := sqlx.Open("postgres", cfg.PostgresAddress)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open Postgres connection")
	}

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "ping failed")
	}

	return &Database{
		Config: cfg,
		Logger: logger,
		TX:     db,
	}, nil
}

func (db *Database) Close() error {
	err := db.TX.(*sqlx.DB).Close()
	return errors.Wrap(err, "cannot close database")
}

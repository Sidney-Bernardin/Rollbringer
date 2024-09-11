package database

import (
	"log/slog"

	"github.com/pkg/errors"

	"rollbringer/internal"
	"rollbringer/internal/config"
	"rollbringer/internal/repositories/databases"
)

type gamesSchema struct {
	*databases.Database[gamesSchema]
}

func New(cfg *config.Config, logger *slog.Logger) (internal.GamesSchema, error) {
	db, err := databases.NewDatabase[gamesSchema](cfg, logger)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create database")
	}

	return &gamesSchema{
		Database: db,
	}, nil
}

package database

import (
	"log/slog"

	"github.com/pkg/errors"

	"rollbringer/internal"
	"rollbringer/internal/config"
	"rollbringer/internal/repositories/databases"
)

type GamesDatabase struct {
	*databases.Database[GamesDatabase]
}

func New(cfg *config.Config, logger *slog.Logger) (internal.GamesDatabase, error) {
	db, err := databases.NewDatabase[GamesDatabase](cfg, logger)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create base database")
	}

	return &GamesDatabase{
		Database: db,
	}, nil
}

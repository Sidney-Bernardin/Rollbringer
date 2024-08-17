package database

import (
	"log/slog"

	"github.com/pkg/errors"

	"rollbringer/internal/config"
	"rollbringer/internal/repositories/databases"
)

type UsersDatabase struct {
	*databases.Database[UsersDatabase]
}

func New(cfg *config.Config, logger *slog.Logger) (*UsersDatabase, error) {

	db, err := databases.NewDatabase[UsersDatabase](cfg, logger)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create base database")
	}

	return &UsersDatabase{
		Database: db,
	}, nil
}

package database

import (
	"log/slog"

	"github.com/pkg/errors"

	"rollbringer/internal"
	"rollbringer/internal/config"
	"rollbringer/internal/repositories/databases"
)

type usersSchema struct {
	*databases.Database[usersSchema]
}

func New(cfg *config.Config, logger *slog.Logger) (internal.UsersSchema, error) {
	db, err := databases.NewDatabase[usersSchema](cfg, logger)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create database")
	}

	return &usersSchema{
		Database: db,
	}, nil
}

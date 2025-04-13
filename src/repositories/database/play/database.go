package play

import (
	"context"
	"rollbringer/src"
	"rollbringer/src/repositories/database"
	"rollbringer/src/services/play"

	"github.com/pkg/errors"
)

type playDatabase struct {
	*database.Database
}

func NewDatabase(ctx context.Context, config *src.Config) (play.Database, error) {
	db, err := database.NewDatabase(ctx, config.PostgresPlayURL)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create database")
	}

	return &playDatabase{db}, nil
}

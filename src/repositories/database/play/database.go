package play

import (
	"embed"

	"rollbringer/src"
	"rollbringer/src/repositories/database"
	"rollbringer/src/services/play"

	"github.com/pkg/errors"
)

//go:embed migrations/*.sql
var migrations embed.FS

type playDatabase struct {
	*database.Database
}

func NewDatabase(config *src.Config) (play.Database, error) {
	database, err := database.NewDatabase(config.PostgresPlayURL, &migrations)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create database")
	}

	return &playDatabase{
		Database: database,
	}, nil
}

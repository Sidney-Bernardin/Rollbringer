package play

import (
	"embed"

	"rollbringer/src"
	"rollbringer/src/domain/play"
	"rollbringer/src/repositories/database"

	"github.com/pkg/errors"
)

//go:embed migrations/*.sql
var migrations embed.FS

type playDatabase struct {
	*database.Database
}

func NewDatabase(config *src.Config) (play.Database, error) {

	database, err := database.NewDatabase(config, &migrations)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create database")
	}

	return &playDatabase{
		Database: database,
	}, nil
}

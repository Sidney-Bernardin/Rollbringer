package play

import (
	"embed"
	"errors"
	"rollbringer/server"
	"rollbringer/server/domain/play"
	"rollbringer/server/repositories/database"
)

//go:embed migrations/*.sql
var migrations embed.FS

type playDatabase struct {
	*database.Database
}

func NewDatabase(config *server.Config) (play.Database, error) {

	database, err := database.NewDatabase(config, &migrations)
	if err != nil {
		return nil, errors.Join(err, errors.New("cannot create database"))
	}

	return &playDatabase{
		Database: database,
	}, nil
}

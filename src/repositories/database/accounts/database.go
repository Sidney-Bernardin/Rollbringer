package accounts

import (
	"embed"

	"github.com/pkg/errors"

	"rollbringer/src"
	"rollbringer/src/repositories/database"
	"rollbringer/src/services/accounts"
)

//go:embed migrations/*.sql
var migrations embed.FS

type accountsDatabase struct {
	*database.Database
}

func NewDatabase(config *src.Config) (accounts.Database, error) {
	database, err := database.NewDatabase(config.PostgresAccountsURL, &migrations)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create database")
	}

	return &accountsDatabase{
		Database: database,
	}, nil
}

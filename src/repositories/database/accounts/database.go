package accounts

import (
	"embed"

	"rollbringer/src"
	"rollbringer/src/domain/accounts"
	"rollbringer/src/repositories/database"

	"github.com/pkg/errors"
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

package accounts

import (
	"context"

	"github.com/pkg/errors"

	"rollbringer/src"
	"rollbringer/src/repositories/database"
	"rollbringer/src/services/accounts"
)

type accountsDatabase struct {
	*database.Database
}

func NewDatabase(ctx context.Context, config *src.Config) (accounts.Database, error) {
	db, err := database.NewDatabase(ctx, config.PostgresAccountsURL)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create database")
	}

	return &accountsDatabase{db}, nil
}

package database

import (
	"context"
	"database/sql"
	"io/fs"
	"log/slog"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/domain/services/accounts"
	"rollbringer/pkg/repositories/database"
)

type accountsDatabaseRepository struct {
	*database.DatabaseRepository

	logger *slog.Logger
}

func NewGamesDatabaseRepository(config *domain.Config, logger *slog.Logger, migrations fs.FS) (service.AccountsDatabaseRepository, error) {
	db, err := database.NewDatabase(config, migrations)
	if err != nil {
		return nil, domain.Wrap(err, "cannot create database", nil)
	}

	return &accountsDatabaseRepository{
		DatabaseRepository: db,
		logger:             logger,
	}, nil
}

func (repo *accountsDatabaseRepository) Transaction(ctx context.Context, txFunc func(service.AccountsDatabaseRepository) error) error {
	tx, err := repo.DB.BeginTxx(ctx, nil)
	if err != nil {
		return domain.Wrap(err, "cannot begin transaction", nil)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			domain.HandleError(ctx, repo.logger, slog.LevelError, domain.Wrap(err, "cannot rollback transaction", nil))
		}
	}()

	err = txFunc(&accountsDatabaseRepository{
		DatabaseRepository: &database.DatabaseRepository{
			TX: tx,
		},
	})

	if err != nil {
		return domain.Wrap(err, "cannot setup transaction", nil)
	}

	err = tx.Commit()
	return domain.Wrap(err, "cannot commit transaction", nil)
}

package database

import (
	"context"
	"database/sql"
	"io/fs"
	"log/slog"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/domain/services/play"
	"rollbringer/pkg/repositories/database"
)

type playDatabaseRepository struct {
	*database.DatabaseRepository

	logger *slog.Logger
}

func NewPlayDatabaseRepository(config *domain.Config, logger *slog.Logger, migrations fs.FS) (service.PlayDatabaseRepository, error) {
	db, err := database.NewDatabase(config, migrations)
	if err != nil {
		return nil, domain.Wrap(err, "cannot create database", nil)
	}

	return &playDatabaseRepository{
		DatabaseRepository: db,
		logger:             logger,
	}, nil
}

func (repo *playDatabaseRepository) Transaction(ctx context.Context, txFunc func(service.PlayDatabaseRepository) error) error {
	tx, err := repo.DB.BeginTxx(ctx, nil)
	if err != nil {
		return domain.Wrap(err, "cannot begin transaction", nil)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			domain.HandleError(ctx, repo.logger, slog.LevelError, domain.Wrap(err, "cannot rollback transaction", nil))
		}
	}()

	err = txFunc(&playDatabaseRepository{
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

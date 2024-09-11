package databases

import (
	"context"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"rollbringer/internal/config"
)

type databaser interface {
	setTX(sqlx.ExtContext)
}

type Database[T databaser] struct {
	Config *config.Config
	Logger *slog.Logger

	TX sqlx.ExtContext
}

func NewDatabase[T databaser](cfg *config.Config, logger *slog.Logger) (*Database[T], error) {
	db, err := sqlx.Open("postgres", cfg.PostgresAddress)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open Postgres connection")
	}

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "ping failed")
	}

	return &Database[T]{
		Config: cfg,
		Logger: logger,
		TX:     db,
	}, nil
}

func (db *Database[T]) setTX(tx sqlx.ExtContext) {
	db.TX = tx
}

func (db *Database[T]) Close() error {
	err := db.TX.(*sqlx.DB).Close()
	return errors.Wrap(err, "cannot close database")
}

func (db *Database[T]) Transaction(ctx context.Context, txFunc func(db *T) error) error {

	tx, err := db.TX.(*sqlx.DB).Beginx()
	if err != nil {
		return errors.Wrap(err, "cannot begin transaction")
	}

	operator := *new(T)
	operator.setTX(tx)

	if err := txFunc(&operator); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return errors.Wrap(rbErr, "cannot rollback failed transaction")
		}

		return errors.Wrap(err, "transaction failed")
	}

	err = tx.Commit()
	return errors.Wrap(err, "cannot commit transaction")
}

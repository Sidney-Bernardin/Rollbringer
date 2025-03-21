package database

import (
	"context"
	"database/sql"
	"rollbringer/src/domain"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type CRUDFunc func(ctx context.Context, view any, q string, args ...any) error

func (db *Database) CRUDInsert(ctx context.Context, view any, query string, args ...any) error {
	if err := sqlx.GetContext(ctx, db.TX, view, query, args...); err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return domain.ErrEntityConflict
			}
		}

		return errors.Wrap(err, "cannot insert row")
	}

	return nil
}

func (db *Database) CRUDGet(ctx context.Context, view any, query string, args ...any) error {
	if err := sqlx.GetContext(ctx, db.TX, view, query, args...); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrEntityNotFound
		}

		return errors.Wrap(err, "cannot get row")
	}

	return nil
}

func (db *Database) CRUDGetMany(ctx context.Context, view any, query string, args ...any) error {
	if err := sqlx.SelectContext(ctx, db.TX, view, query, args...); err != nil {
		return errors.Wrap(err, "cannot get rows")
	}
	return nil
}

func (db *Database) CRUDUpdate(ctx context.Context, view any, query string, args ...any) error {

	result, err := db.TX.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "cannot update row")
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "cannot get affected rows")
	}

	if affected <= 0 {
		return domain.ErrNoEntitiesEffected
	}

	return nil
}

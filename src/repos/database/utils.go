package database

import (
	"context"
	"database/sql"
	"rollbringer/src/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
)

type Row[M any] interface {
	Model() M
}

type Tx interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func Insert(ctx context.Context, tx Tx, query string, args ...any) error {
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return domain.ErrEntityConflict
			}
		}

		return errors.Wrap(err, "cannot insert")
	}

	return nil
}

func Get[R Row[M], M any](ctx context.Context, tx Tx, query string, args ...any) (row *R, model M, err error) {
	q, err := tx.Query(ctx, query, args...)
	if err != nil {
		return row, model, errors.Wrap(err, "cannot query")
	}

	row, err = pgx.CollectExactlyOneRow(q, pgx.RowToAddrOfStructByNameLax[R])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return row, model, domain.ErrEntityNotFound
		}

		return row, model, errors.Wrap(err, "cannot get row")
	}

	model = (*row).Model()
	return row, model, nil
}

func Gets[R interface{ Model() M }, M any](ctx context.Context, tx Tx, query string, args ...any) ([]*R, []M, error) {
	q, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot query")
	}

	rows, err := pgx.CollectRows(q, pgx.RowToAddrOfStructByNameLax[R])
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot get rows")
	}

	models := make([]M, len(rows))
	for i, row := range rows {
		models[i] = (*row).Model()
	}

	return rows, models, nil
}

func Update(ctx context.Context, tx Tx, query string, args ...any) error {
	result, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "cannot update")
	}

	if result.RowsAffected() < 1 {
		return domain.ErrNoEntitiesEffected
	}

	return nil
}

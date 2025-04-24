package database

import (
	"context"
	"database/sql"
	"embed"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"

	"rollbringer/src"
)

//go:embed Migrations/*.sql
var migrations embed.FS

type Database struct {
	Pool *pgxpool.Pool
	Tx   Transaction
}

type Transaction interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func NewDatabase(ctx context.Context, dbURL string) (*Database, error) {

	// Connect to Postgres.
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to Postgres server")
	}

	// Create migration source.
	migrationSrc, err := iofs.New(migrations, "Migrations")
	if err != nil {
		return nil, errors.Wrap(err, "cannot create migration source")
	}

	// Create migrate instance.
	migrator, err := migrate.NewWithSourceInstance("iofs", migrationSrc, dbURL)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create migrate instance")
	}

	// Migrate database.
	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, errors.Wrap(err, "cannot migrate")
	}

	return &Database{pool, pool}, nil
}

func (db *Database) Close() {
	db.Pool.Reset()
}

func (db *Database) Transaction(ctx context.Context, txFunc func(tx pgx.Tx) error) error {

	// Begin transaction.
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot begin transaction")
	}
	defer tx.Rollback(ctx)

	// Do transaction callback.
	if err := txFunc(tx); err != nil {
		return errors.Wrap(err, "transaction failed")
	}

	// Commit transaction.
	err = tx.Commit(ctx)
	return errors.Wrap(err, "cannot commit transaction")
}

func Insert(ctx context.Context, tx Transaction, query string, args ...any) error {
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return src.ErrEntityConflict
			}
		}

		return errors.Wrap(err, "cannot insert")
	}

	return nil
}

func Get[T any](ctx context.Context, tx Transaction, query string, args ...any) (*T, error) {
	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "cannot query")
	}

	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[T])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, src.ErrEntityNotFound
		}

		return nil, errors.Wrap(err, "cannot get row")
	}

	return &res, nil
}

func Gets[T any](ctx context.Context, tx Transaction, query string, args ...any) ([]*T, error) {
	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "cannot query")
	}

	res, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[T])
	return res, errors.Wrap(err, "cannot get rows")
}

func Update(ctx context.Context, tx Transaction, query string, args ...any) error {
	result, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "cannot update")
	}

	if result.RowsAffected() < 1 {
		return src.ErrNoEntitiesEffected
	}

	return nil
}

func Domains[R interface{ Domain() D }, D any](rows []R) []D {
	ret := make([]D, len(rows))
	for i, row := range rows {
		ret[i] = row.Domain()
	}
	return ret
}

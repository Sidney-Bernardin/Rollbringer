package database

import (
	"context"
	"database/sql"
	"rollbringer/pkg/domain"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/lib/pq/hstore"
	"github.com/pkg/errors"
)

type conn interface {
	Exec(sql string, arguments ...any) (sql.Result, error)
	Query(sql string, args ...any) (*sql.Rows, error)
	QueryRow(sql string, args ...any) *sql.Row
}

type Database struct {
	conn conn
}

// New returns a new Database that connects to a Postgres server.
func New(addr string) (*Database, error) {

	// Connect to the Postgres server.
	conn, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to postgres server")
	}

	// Ping the Postgres server.
	if err := conn.Ping(); err != nil {
		return nil, errors.Wrap(err, "cannot ping postgres server")
	}

	return &Database{conn}, nil
}

func (db *Database) Transaction(ctx context.Context, txFunc func(db *Database) error) error {

	tx, err := db.conn.(*sql.DB).Begin()
	if err != nil {
		return errors.Wrap(err, "cannot begin transaction")
	}

	if err := txFunc(&Database{tx}); err != nil {

		if rbErr := tx.Rollback(); rbErr != nil {
			return errors.Wrapf(rbErr, `transaction error (%v) caused rollback error (%v)`, err, rbErr)
		}

		return errors.Wrap(err, "transaction failed")
	}

	return errors.Wrap(tx.Commit(), "cannot commit transaction")
}

func (db *Database) parseUUIDs(ids ...*string) {
	for _, id := range ids {
		parsed, _ := uuid.Parse(*id)
		*id = parsed.String()
	}
}

// decodeHstorePages copies the hstore pages into the PDF.
func decodeHstorePages(hstorePages []hstore.Hstore, pdf *domain.PDF) {
	pdf.Pages = make([]map[string]string, len(hstorePages))
	for i, page := range hstorePages {
		pdf.Pages[i] = map[string]string{}
		for k, v := range page.Map {
			pdf.Pages[i][k] = v.String
		}
	}
}

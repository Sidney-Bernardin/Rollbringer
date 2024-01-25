package database

import (
	"context"
	"fmt"
	"rollbringer/pkg/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

// InsertPDF inserts a new PDF for the owner.
func (db *Database) InsertPDF(ctx context.Context, ownerID uuid.UUID, schema string) (uuid.UUID, string, error) {

	pdfID := uuid.New()
	name := fmt.Sprintf("Unnamed %s", schema)

	// Insert a new PDF for the owner.
	_, err := db.conn.Exec(ctx,
		`INSERT INTO pdfs (id, owner_id, name, schema, content) VALUES ($1, $2, $3, $4, $5)`,
		pdfID, ownerID, name, schema, []byte(``))

	if err != nil {
		return uuid.Nil, "", errors.Wrap(err, "cannot insert pdf")
	}

	return pdfID, name, nil
}

// GetPDFs return the PDFs with the owner-ID from the database.
func (db *Database) GetPDFs(ctx context.Context, ownerID uuid.UUID) ([]*domain.PDF, error) {

	// Get the pdfs with the owner-ID.
	rows, err := db.conn.Query(ctx, `SELECT * FROM pdfs WHERE owner_id = $1`, ownerID)
	if err != nil {
		return nil, errors.Wrap(err, "cannot select pdfs")
	}
	defer rows.Close()

	// Scan into a slice of PDF models.
	pdfs, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[domain.PDF])
	return pdfs, errors.Wrap(err, "cannot scan pdfs")
}

// DeletePDF deletes the pdf with the pdf-ID and owner-ID from the database.
// If the pdf doesn't exist, returns domain.ErrPlayMaterialNotFound.
func (db *Database) DeletePDF(ctx context.Context, pdfID, ownerID uuid.UUID) error {

	// Delete the pdf with the pdf-ID and owner-ID.
	cmdTag, err := db.conn.Exec(ctx,
		`DELETE FROM pdfs WHERE id = $1 AND owner_id = $2`,
		pdfID, ownerID)

	if err != nil {
		return errors.Wrap(err, "cannot delete pdf")
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.ErrGameNotFound
	}

	return nil
}


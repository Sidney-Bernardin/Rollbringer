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
func (db *Database) InsertPDF(ctx context.Context, ownerID string, schema string) (string, string, error) {

	ownerUUID, _ := uuid.Parse(ownerID)
	pdfID := uuid.New().String()
	name := fmt.Sprintf("Unnamed %s", schema)

	// Insert a new PDF for the owner.
	_, err := db.conn.Exec(ctx,
		`INSERT INTO pdfs (id, owner_id, name, schema, main_page, info_page, spells_page) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		pdfID, ownerUUID, name, schema, []byte(``), []byte(``), []byte(``))

	if err != nil {
		return "", "", errors.Wrap(err, "cannot insert pdf")
	}

	return pdfID, name, nil
}

// GetPDF returns the PDF with the PDF-ID from the database. If the game
// doesn't exist, returns domain.ErrPlayMaterialNotFound.
func (db *Database) GetPDF(ctx context.Context, pdfID, ownerID string) (*domain.PDF, error) {

	pdfUUID, _ := uuid.Parse(pdfID)
	ownerUUID, _ := uuid.Parse(ownerID)

	// Get the PDF with the PDF-ID.
	rows, err := db.conn.Query(ctx, `SELECT * FROM pdfs WHERE id = $1 AND owner_id = $2`, pdfUUID, ownerUUID)
	if err != nil {
		return nil, errors.Wrap(err, "cannot select pdf")
	}
	defer rows.Close()

	// Scan into a PDF model.
	pdf, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[domain.PDF])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPlayMaterialNotFound
		}

		return nil, errors.Wrap(err, "cannot scan pdf")
	}

	return pdf, nil
}

// GetPDFs return the PDFs with the owner-ID from the database.
func (db *Database) GetPDFs(ctx context.Context, ownerID string) ([]*domain.PDF, error) {

	ownerUUID, _ := uuid.Parse(ownerID)

	// Get the PDFs with the owner-ID.
	rows, err := db.conn.Query(ctx, `SELECT * FROM pdfs WHERE owner_id = $1`, ownerUUID)
	if err != nil {
		return nil, errors.Wrap(err, "cannot select pdfs")
	}

	defer rows.Close()

	// Scan into a slice of PDF models.
	pdfs, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[domain.PDF])
	return pdfs, errors.Wrap(err, "cannot scan pdfs")
}

// UpdatePDF updates the PDF with the pdf-ID and owner-ID in the database.
func (db *Database) UpdatePDF(ctx context.Context, page int, pdfID, ownerID string, content []byte) error {

	//	pdfUUID, _ := uuid.Parse(pdfID)
	//	ownerUUID, _ := uuid.Parse(ownerID)

	//	if cmdTag.RowsAffected() == 0 {
	//		return domain.ErrPlayMaterialNotFound
	//	}

	//	return errors.Wrap(err, "cannot update pdf")

	return nil
}

// DeletePDF deletes the pdf with the pdf-ID and owner-ID from the database.
// If the pdf doesn't exist, returns domain.ErrPlayMaterialNotFound.
func (db *Database) DeletePDF(ctx context.Context, pdfID, ownerID string) error {

	pdfUUID, _ := uuid.Parse(pdfID)
	ownerUUID, _ := uuid.Parse(ownerID)

	// Delete the pdf with the pdf-ID and owner-ID.
	cmdTag, err := db.conn.Exec(ctx,
		`DELETE FROM pdfs WHERE id = $1 AND owner_id = $2`,
		pdfUUID, ownerUUID)

	if err != nil {
		return errors.Wrap(err, "cannot delete pdf")
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.ErrPlayMaterialNotFound
	}

	return nil
}

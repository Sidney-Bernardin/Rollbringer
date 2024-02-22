package database

import (
	"context"
	"encoding/json"
	"rollbringer/pkg/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

// InsertPDF inserts the PDF.
func (db *Database) InsertPDF(ctx context.Context, pdf *domain.PDF) error {

	ownerUUID, _ := uuid.Parse(pdf.OwnerID)
	pdf.ID = uuid.New().String()

	// Insert a new PDF.
	_, err := db.conn.Exec(ctx,
		`INSERT INTO pdfs (id, owner_id, name, schema, pages) VALUES ($1, $2, $3, $4, $5)`,
		pdf.ID, ownerUUID, pdf.Name, pdf.Schema, []string{"{}", "{}", "{}"})

	return errors.Wrap(err, "cannot insert pdf")
}

// GetPDF returns the PDF with the PDF-ID. If the PDF doesn't exist,
// returns domain.ErrPlayMaterialNotFound.
func (db *Database) GetPDF(ctx context.Context, pdfID string) (*domain.PDF, error) {

	pdfUUID, _ := uuid.Parse(pdfID)

	// Get the PDF with the PDF-ID.
	rows, err := db.conn.Query(ctx, `SELECT * FROM pdfs WHERE id = $1`, pdfUUID)
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

// GetPDFs return the PDFs with the owner-ID.
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

// UpdatePDFPage updates the PDF with the PDF-ID. If the PDF doesn't exist,
// returns domain.ErrPlayMaterialNotFound.
func (db *Database) UpdatePDFPage(ctx context.Context, pdfID string, pageNum int, pdfPage any) error {

	pdfUUID, _ := uuid.Parse(pdfID)

	// Encode the PDF page.
	pdfFieldsJSON, err := json.Marshal(pdfPage)
	if err != nil {
		return errors.Wrap(err, "cannot encode pdf page")
	}

	// Update the PDF with the PDF-ID.
	cmdTag, err := db.conn.Exec(ctx,
		`UPDATE pdfs SET pages[$1] = $2 WHERE id = $3`,
		pageNum-1, string(pdfFieldsJSON), pdfUUID)

	if err != nil {
		return errors.Wrap(err, "cannot update pdf")
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.ErrPlayMaterialNotFound
	}

	return nil
}

// DeletePDF deletes the PDF with the PDF-ID and owner-ID. If the PDF doesn't
// exist, returns domain.ErrPlayMaterialNotFound.
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

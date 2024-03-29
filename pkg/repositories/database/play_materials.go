package database

import (
	"context"
	"database/sql"
	"rollbringer/pkg/domain"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/lib/pq/hstore"
	"github.com/pkg/errors"
)

type pdfModel struct {
	domain.PDF
	hstorePages []hstore.Hstore
}

// InsertPDF inserts the PDF.
func (db *Database) InsertPDF(ctx context.Context, pdf *domain.PDF) error {
	pdf.ID = uuid.New().String()

	// Create an HStore copy of the PDF's pages.
	hstorePages := make([]hstore.Hstore, len(pdf.Pages))
	for i := range hstorePages {
		hstorePages[i].Map = map[string]sql.NullString{}
	}

	// Insert a new PDF.
	_, err := db.conn.Exec(
		`INSERT INTO pdfs (id, owner_id, game_id, name, schema, pages) 
			VALUES ($1, $2, $3, $4, $5, $6)`,
		pdf.ID, pdf.OwnerID, pdf.GameID, pdf.Name, pdf.Schema, pq.Array(hstorePages))

	return errors.Wrap(err, "cannot insert pdf")
}

// GetPDF returns the PDF with the PDF-ID.
func (db *Database) GetPDF(ctx context.Context, pdfID string) (*domain.PDF, error) {

	var (
		pdf         domain.PDF
		hstorePages []hstore.Hstore
	)

	// Get the PDF with the PDF-ID.
	err := db.conn.QueryRow(
		`SELECT id, owner_id, game_id, name, schema, pages FROM pdfs WHERE id = $1`, pdfID).
		Scan(&pdf.ID, &pdf.OwnerID, &pdf.GameID, &pdf.Name, &pdf.Schema, pq.Array(&hstorePages))

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &domain.ProblemDetail{
				Type:   domain.PDTypePDFNotFound,
				Detail: "No PDF with the given PDF-ID exists was found.",
			}
		}

		return nil, errors.Wrap(err, "cannot select pdf")
	}

	// Set the PDF's pages to the HStore pages.
	for i := range hstorePages {
		pdf.Pages = append(pdf.Pages, hstoreToMap(hstorePages[i]))
	}

	return &pdf, nil
}

// GetPDFs return the PDFs with the owner-ID.
func (db *Database) GetPDFs(ctx context.Context, ownerID string) ([]*domain.PDF, error) {

	// Get the PDFs with the owner-ID.
	rows, err := db.conn.Query(
		`SELECT id, owner_id, game_id, name, schema, pages FROM pdfs WHERE owner_id = $1`,
		ownerID)

	if err != nil {
		return nil, errors.Wrap(err, "cannot select pdfs")
	}
	defer rows.Close()

	// Scan the rows into a slice of PDFs.
	pdfs := []*domain.PDF{}
	for rows.Next() {

		var (
			pdf         domain.PDF
			hstorePages []hstore.Hstore
		)

		// Scan the row.
		err := rows.Scan(&pdf.ID, &pdf.OwnerID, &pdf.GameID, &pdf.Name, &pdf.Schema, pq.Array(&hstorePages))
		if err != nil {
			return nil, errors.Wrap(err, "cannot scan pdf")
		}

		// Set the PDF's pages to the HStore pages.
		for i := range hstorePages {
			pdf.Pages = append(pdf.Pages, hstoreToMap(hstorePages[i]))
		}

		pdfs = append(pdfs, &pdf)
	}

	return pdfs, nil
}

// UpdatePDFField updates the page field of the PDF with the PDF-ID.
func (db *Database) UpdatePDFField(ctx context.Context, pdfID string, pageIdx int, fieldName, fieldValue string) error {

	// Update the page field of the PDF with the PDF-ID.
	result, err := db.conn.Exec(
		`UPDATE pdfs SET pages[$1] = pages[$1] || hstore($2, $3) WHERE id = $4`,
		pageIdx+1, fieldName, fieldValue, pdfID)

	if err != nil {
		return errors.Wrap(err, "cannot update pdf")
	}

	// Get the number of rows affected.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "cannot get rows affected")
	}

	if rowsAffected == 0 {
		return &domain.ProblemDetail{
			Type:   domain.PDTypePDFNotFound,
			Detail: "No PDF with the given PDF-ID was found.",
		}
	}

	return nil
}

// DeletePDF deletes the PDF with the PDF-ID and owner-ID.
func (db *Database) DeletePDF(ctx context.Context, pdfID, ownerID string) error {

	// Delete the pdf with the pdf-ID and owner-ID.
	result, err := db.conn.Exec(
		`DELETE FROM pdfs WHERE id = $1 AND owner_id = $2`,
		pdfID, ownerID)

	if err != nil {
		return errors.Wrap(err, "cannot delete pdf")
	}

	// Get the number of rows affected.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "cannot get rows affected")
	}

	if rowsAffected == 0 {
		return &domain.ProblemDetail{
			Type:   domain.PDTypePDFNotFound,
			Detail: "No PDF with the given PDF-ID and owner-ID was found.",
		}
	}

	return nil
}

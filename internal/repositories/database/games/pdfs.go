package games

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/lib/pq/hstore"
	"github.com/pkg/errors"

	"rollbringer/internal"
	"rollbringer/internal/repositories/database"
)

func pdfColumns(views map[string]internal.PDFView) (columns string) {
	switch views["pdf"] {
	case internal.PDFViewPDFAll:
		columns += `pdfs.id, pdfs.owner_id, pdfs.game_id, pdfs.name, pdfs.schema`
	}

	switch views["owner"] {
	case internal.PDFViewOwnerInfo:
		columns += `users.id AS "owner.id", COALESCE(users.username, '') AS "owner.username"`
	}

	switch views["game"] {
	case internal.PDFViewGameInfo:
		columns += `games.id AS "game.id", COALESCE(games.name, '') AS "game.name"`
	}

	return columns
}

func pdfJoins(views map[string]internal.PDFView) (joins string) {
	if _, ok := views["owner"]; ok {
		joins += `LEFT JOIN users ON users.id = pdfs.owner_id`
	}

	if _, ok := views["host"]; ok {
		joins += `LEFT JOIN games ON games.id = pdfs.game_id`
	}

	return joins
}

func (db *gamesSchema) PDFInsert(ctx context.Context, pdf *internal.PDF) error {
	hstorePages := make([]hstore.Hstore, len(pdf.Pages))
	for i, page := range pdf.Pages {
		hstorePage := map[string]sql.NullString{}
		if page != nil {
			for k, v := range page {
				hstorePage[k] = sql.NullString{String: v, Valid: true}
			}
		}
		hstorePages[i].Map = hstorePage
	}

	err := db.TX.QueryRowxContext(ctx,
		`INSERT INTO pdfs (id, owner_id, game_id, name, schema, pages)
			VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		uuid.New(), pdf.OwnerID, pdf.GameID, pdf.Name, pdf.Schema, pq.Array(hstorePages),
	).Scan(&pdf.ID)

	return errors.Wrap(err, "cannot insert PDF")
}

func (db *gamesSchema) PDFGet(ctx context.Context, pdfID uuid.UUID, view map[string]internal.PDFView) (*internal.PDF, error) {

	var pdf database.PDF
	query := fmt.Sprintf(`SELECT %s FROM pdfs %s WHERE pdfs.id = $1`, pdfColumns(view), pdfJoins(view))
	if err := sqlx.GetContext(ctx, db.TX, &pdf, query, pdfID); err != nil {
		if err == sql.ErrNoRows {
			return nil, internal.NewProblemDetail(ctx, internal.PDOpts{
				Type:   internal.PDTypePDFNotFound,
				Detail: "Can't find a PDF with the given pdf_id.",
				Extra: map[string]any{
					"pdf_id": pdfID,
				},
			})
		}

		return nil, errors.Wrap(err, "cannot select PDF")
	}

	return pdf.Internalized(), nil
}

func (db *gamesSchema) PDFGetPage(ctx context.Context, pdfID uuid.UUID, pageIdx int) (map[string]string, error) {

	var page hstore.Hstore
	err := db.TX.QueryRowxContext(ctx,
		`SELECT pages[$1] FROM pdfs WHERE id = $2`,
		pageIdx+1, pdfID,
	).Scan(&page)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, internal.NewProblemDetail(ctx, internal.PDOpts{
				Type:   internal.PDTypePDFNotFound,
				Detail: "Can't find a PDF with the given pdf_id.",
				Extra: map[string]any{
					"pdf_id": pdfID,
				},
			})
		}

		return nil, errors.Wrap(err, "cannot get PDF page")
	}

	// Convert the page to a map without null-strings.
	ret := make(map[string]string, len(page.Map))
	for k, v := range page.Map {
		ret[k] = v.String
	}

	return ret, nil
}

func (db *gamesSchema) PDFsGetForOwner(ctx context.Context, ownerID uuid.UUID, views map[string]internal.PDFView) ([]*internal.PDF, error) {

	var pdfs []*database.PDF
	query := fmt.Sprintf(`SELECT %s FROM pdfs %s WHERE pdfs.owner_id = $1`, pdfColumns(views), pdfJoins(views))
	if err := sqlx.SelectContext(ctx, db.TX, &pdfs, query, ownerID); err != nil {
		return nil, errors.Wrap(err, "cannot select PDFs")
	}

	// Internalize each pdf.
	ret := make([]*internal.PDF, len(pdfs))
	for i, m := range pdfs {
		ret[i] = m.Internalized()
	}

	return ret, nil
}

func (db *gamesSchema) PDFsGetForGame(ctx context.Context, gameID uuid.UUID, views map[string]internal.PDFView) ([]*internal.PDF, error) {

	var pdfs []*database.PDF
	query := fmt.Sprintf(`SELECT %s FROM pdfs %s WHERE pdfs.game_id = $1`, pdfColumns(views), pdfJoins(views))
	if err := sqlx.SelectContext(ctx, db.TX, &pdfs, query, gameID); err != nil {
		return nil, errors.Wrap(err, "cannot select PDFs")
	}

	// Internalize each pdf.
	ret := make([]*internal.PDF, len(pdfs))
	for i, m := range pdfs {
		ret[i] = m.Internalized()
	}

	return ret, nil
}

func (db *gamesSchema) PDFUpdatePage(ctx context.Context, pdfID uuid.UUID, pageIdx int, fieldName, fieldValue string) error {
	_, err := db.TX.ExecContext(ctx,
		`UPDATE pdfs SET pages[$1] = pages[$1] || hstore($2, $3) WHERE id = $4`,
		pageIdx+1, fieldName, fieldValue, pdfID)

	return errors.Wrap(err, "cannot update PDF page")
}

func (db *gamesSchema) PDFDelete(ctx context.Context, pdfID, ownerID uuid.UUID) error {
	_, err := db.TX.ExecContext(ctx,
		`DELETE FROM pdfs WHERE id = $1 AND owner_id = $2`,
		pdfID, ownerID)

	return errors.Wrap(err, "cannot delete PDF")
}
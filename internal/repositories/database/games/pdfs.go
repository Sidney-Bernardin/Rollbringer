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

func pdfColumns(view internal.PDFView) string {
	switch view {
	case internal.PDFViewListItem:
		return `pdfs.id, pdfs.owner_id, pdfs.game_id, pdfs.name, pdfs.schema,` +
			`games.id AS "game.id",` +
			`COALESCE(games.name, '') AS "game.name",` +
			`users.id AS "owner.id",` +
			`users.username AS "owner.username",` +
			`users.google_id AS "owner.google_id"`

	default:
		return `pdfs.id, pdfs.owner_id, pdfs.game_id, pdfs.name, pdfs.schema`
	}
}

func pdfJoins(view internal.PDFView) string {
	switch view {
	case internal.PDFViewListItem:
		return `LEFT JOIN games.games ON games.id = pdfs.game_id LEFT JOIN users.users ON users.id = pdfs.owner_id`
	default:
		return ``
	}
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
		`INSERT INTO games.pdfs (id, owner_id, game_id, name, schema, pages)
			VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		uuid.New(), pdf.OwnerID, pdf.GameID, pdf.Name, pdf.Schema, pq.Array(hstorePages),
	).Scan(&pdf.ID)

	return errors.Wrap(err, "cannot insert PDF")
}

func (db *gamesSchema) PDFGet(ctx context.Context, pdfID uuid.UUID, view internal.PDFView) (*internal.PDF, error) {

	var pdf database.PDF
	query := fmt.Sprintf(`SELECT %s FROM games.pdfs %s WHERE pdfs.id = $1`, pdfColumns(view), pdfJoins(view))
	if err := sqlx.GetContext(ctx, db.TX, &pdf, query, pdfID); err != nil {
		if err == sql.ErrNoRows {
			return nil, internal.NewProblemDetail(ctx, internal.PDOpts{
				Type:   internal.PDTypePDFNotFound,
				Detail: "Cannot find a PDF with the given pdf_id.",
				Extra: map[string]any{
					"pdf_id": pdfID,
				},
			})
		}

		return nil, errors.Wrap(err, "cannot select PDF")
	}

	return pdf.Internalized(), nil
}

func (db *gamesSchema) PDFGetPage(ctx context.Context, pdfID uuid.UUID, pageNum int) (map[string]string, error) {

	var page hstore.Hstore
	err := db.TX.QueryRowxContext(ctx,
		`SELECT pages[$1] FROM games.pdfs WHERE id = $2`,
		pageNum, pdfID,
	).Scan(&page)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, internal.NewProblemDetail(ctx, internal.PDOpts{
				Type:   internal.PDTypePDFNotFound,
				Detail: "Cannot find a PDF with the given pdf_id.",
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

func (db *gamesSchema) PDFsGetByOwner(ctx context.Context, ownerID uuid.UUID, view internal.PDFView) ([]*internal.PDF, error) {

	var pdfs []*database.PDF
	query := fmt.Sprintf(`SELECT %s FROM games.pdfs %s WHERE pdfs.owner_id = $1`, pdfColumns(view), pdfJoins(view))
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

func (db *gamesSchema) PDFsGetByGame(ctx context.Context, gameID uuid.UUID, view internal.PDFView) ([]*internal.PDF, error) {

	var pdfs []*database.PDF
	query := fmt.Sprintf(`SELECT %s FROM games.pdfs %s WHERE pdfs.game_id = $1`, pdfColumns(view), pdfJoins(view))
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

func (db *gamesSchema) PDFUpdate(ctx context.Context, session *internal.Session, pdf *internal.PDF) error {
	var sets string

	if pdf.Name != "" {
		sets += `SET name = :name`
	}

	if sets == "" {
		return nil
	}

	query := fmt.Sprintf(`UPDATE games.pdfs %s WHERE id = :id AND owner_id = :owner_id`, sets)
	result, err := sqlx.NamedExecContext(ctx, db.TX, query, map[string]any{
		"name":     pdf.Name,
		"id":       pdf.ID,
		"owner_id": session.UserID,
	})

	if err != nil {
		return errors.Wrap(err, "cannot update PDF")
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "cannot check affected rows")
	}

	if affected <= 0 {
		return internal.NewProblemDetail(ctx, internal.PDOpts{
			Type:   internal.PDTypePDFNotFound,
			Detail: "Cannot find a PDF with the given pdf_id.",
			Extra: map[string]any{
				"pdf_id": pdf.ID,
			},
		})
	}

	return nil
}

func (db *gamesSchema) PDFUpdatePage(ctx context.Context, pdfID uuid.UUID, pageNum int, fieldName, fieldValue string) error {
	result, err := db.TX.ExecContext(ctx,
		`UPDATE games.pdfs SET pages[$1] = pages[$1] || hstore($2, $3) WHERE id = $4`,
		pageNum, fieldName, fieldValue, pdfID)

	if err != nil {
		return errors.Wrap(err, "cannot update PDF page")
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "cannot check affected rows")
	}

	if affected <= 0 {
		return internal.NewProblemDetail(ctx, internal.PDOpts{
			Type:   internal.PDTypePDFNotFound,
			Detail: "Cannot find a PDF with the given pdf_id.",
			Extra: map[string]any{
				"pdf_id": pdfID,
			},
		})
	}

	return nil
}

func (db *gamesSchema) PDFDelete(ctx context.Context, pdfID, ownerID uuid.UUID) error {
	_, err := db.TX.ExecContext(ctx,
		`DELETE FROM games.pdfs WHERE id = $1 AND owner_id = $2`,
		pdfID, ownerID)

	return errors.Wrap(err, "cannot delete PDF")
}

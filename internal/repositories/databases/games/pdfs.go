package database

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
)

var pdfViewColumns = map[internal.PDFView]string{
	internal.PDFViewAll:          `pdfs.id, pdfs.owner_id, pdfs.game_id, pdfs.name, pdfs.schema`,
	internal.PDFViewAll_GameInfo: `pdfs.id, pdfs.owner_id, pdfs.game_id, pdfs.name, pdfs.schema, games.id AS "game.id", COALESCE(games.name, '') AS "game.name"`,
}

type dbPDF struct {
	ID uuid.UUID `db:"id"`

	OwnerID uuid.UUID `db:"owner_id"`

	GameID *uuid.UUID `db:"game_id"`
	Game   *dbGame    `db:"game"`

	Name   string          `db:"name"`
	Schema string          `db:"schema"`
	Pages  pq.GenericArray `db:"pages"`
}

func (pdf *dbPDF) internalized() *internal.PDF {
	if pdf != nil {
		return &internal.PDF{
			ID:      pdf.ID,
			OwnerID: pdf.OwnerID,
			GameID:  pdf.GameID,
			Game:    pdf.Game.internalized(),
			Name:    pdf.Name,
			Schema:  pdf.Schema,
		}
	}
	return nil
}

func (db *GamesDatabase) PDFInsert(ctx context.Context, pdf *internal.PDF, pageCount int) error {
	hstorePages := make([]hstore.Hstore, pageCount)
	for i := range pageCount {
		hstorePages[i].Map = map[string]sql.NullString{}
	}

	err := db.TX.QueryRowxContext(ctx,
		`INSERT INTO pdfs (id, owner_id, game_id, name, schema, pages)
			VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		uuid.New(), pdf.OwnerID, pdf.GameID, pdf.Name, pdf.Schema, pq.Array(hstorePages),
	).Scan(&pdf.ID)

	return errors.Wrap(err, "cannot insert PDF")
}

func (db *GamesDatabase) PDFsGetForOwner(ctx context.Context, ownerID uuid.UUID, view internal.PDFView) ([]*internal.PDF, error) {

	var joins string
	if view == internal.PDFViewAll_GameInfo {
		joins = `LEFT JOIN games ON games.id = pdfs.game_id`
	}

	columns, ok := pdfViewColumns[view]
	if !ok {
		return nil, fmt.Errorf("bad PDF view %d", view)
	}
	query := fmt.Sprintf(`SELECT %s FROM pdfs %s WHERE pdfs.owner_id = $1`, columns, joins)

	var pdfs []*dbPDF
	if err := sqlx.SelectContext(ctx, db.TX, &pdfs, query, ownerID); err != nil {
		return nil, errors.Wrap(err, "cannot select PDFs")
	}

	// Internalize each pdf.
	ret := make([]*internal.PDF, len(pdfs))
	for i, m := range pdfs {
		ret[i] = m.internalized()
	}

	return ret, nil
}

func (db *GamesDatabase) PDFsGetForGame(ctx context.Context, gameID uuid.UUID, view internal.PDFView) ([]*internal.PDF, error) {

	var joins string
	if view == internal.PDFViewAll_GameInfo {
		joins = `LEFT JOIN games ON games.id = pdfs.game_id`
	}

	columns, ok := pdfViewColumns[view]
	if !ok {
		return nil, fmt.Errorf("bad PDF view %d", view)
	}
	query := fmt.Sprintf(`SELECT %s FROM pdfs %s WHERE pdfs.game_id = $1`, columns, joins)

	var pdfs []*dbPDF
	if err := sqlx.SelectContext(ctx, db.TX, &pdfs, query, gameID); err != nil {
		return nil, errors.Wrap(err, "cannot select PDFs")
	}

	// Internalize each pdf.
	ret := make([]*internal.PDF, len(pdfs))
	for i, m := range pdfs {
		ret[i] = m.internalized()
	}

	return ret, nil
}

func (db *GamesDatabase) PDFGet(ctx context.Context, pdfID uuid.UUID, view internal.PDFView) (*internal.PDF, error) {

	var joins string
	if view == internal.PDFViewAll_GameInfo {
		joins = `LEFT JOIN games ON games.id = pdfs.game_id`
	}

	columns, ok := pdfViewColumns[view]
	if !ok {
		return nil, fmt.Errorf("bad PDF view %d", view)
	}
	query := fmt.Sprintf(`SELECT %s FROM pdfs %s WHERE pdfs.id = $1`, columns, joins)

	var pdf dbPDF
	if err := sqlx.GetContext(ctx, db.TX, &pdf, query, pdfID); err != nil {
		if err == sql.ErrNoRows {
			return nil, internal.NewProblemDetail(ctx, &internal.PDOptions{
				Type:   internal.PDTypePDFNotFound,
				Detail: "Can't find a PDF with the given PDF-ID.",
				Extra: map[string]any{
					"pdf_id": pdfID,
				},
			})
		}

		return nil, errors.Wrap(err, "cannot select PDF")
	}

	return pdf.internalized(), nil
}

func (db *GamesDatabase) PDFGetPage(ctx context.Context, pdfID uuid.UUID, pageIdx int) (map[string]string, error) {

	var page hstore.Hstore
	err := db.TX.QueryRowxContext(ctx,
		`SELECT pages[$1] FROM pdfs WHERE id = $2`,
		pageIdx+1, pdfID,
	).Scan(&page)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, internal.NewProblemDetail(ctx, &internal.PDOptions{
				Type:   internal.PDTypePDFNotFound,
				Detail: "Can't find a PDF with the given PDF-ID.",
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

func (db *GamesDatabase) PDFUpdatePage(ctx context.Context, pdfID uuid.UUID, pageIdx int, fieldName, fieldValue string) error {
	_, err := db.TX.ExecContext(ctx,
		`UPDATE pdfs SET pages[$1] = pages[$1] || hstore($2, $3) WHERE id = $4`,
		pageIdx+1, fieldName, fieldValue, pdfID)

	return errors.Wrap(err, "cannot update PDF page")
}

func (db *GamesDatabase) PDFDelete(ctx context.Context, pdfID, ownerID uuid.UUID) error {
	_, err := db.TX.ExecContext(ctx,
		`DELETE FROM pdfs WHERE id = $1 AND owner_id = $2`,
		pdfID, ownerID)

	return errors.Wrap(err, "cannot delete PDF")
}

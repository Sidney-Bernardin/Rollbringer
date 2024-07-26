package database

import (
	"context"
	"database/sql"
	"fmt"
	"rollbringer/pkg/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/lib/pq/hstore"
	"github.com/pkg/errors"
)

var pdfViewColumns = map[domain.PDFView]string{
	domain.PDFViewAll:                    `pdfs.id, pdfs.owner_id, pdfs.game_id, pdfs.name, pdfs.schema`,
	domain.PDFViewAll_OwnerInfo:          `pdfs.id, pdfs.owner_id, pdfs.game_id, pdfs.name, pdfs.schema, users.id AS "owner.id", COALESCE(users.username, '') AS "owner.username"`,
	domain.PDFViewAll_GameInfo:           `pdfs.id, pdfs.owner_id, pdfs.game_id, pdfs.name, pdfs.schema, games.id AS "game.id", COALESCE(games.name, '') AS "game.name"`,
	domain.PDFViewAll_OwnerInfo_GameInfo: `pdfs.id, pdfs.owner_id, pdfs.game_id, pdfs.name, pdfs.schema, users.id AS "owner.id", COALESCE(users.username, '') AS "owner.username", games.id AS "game.id", COALESCE(games.name, '') AS "game.name"`,
}

type pdfModel struct {
	ID uuid.UUID `db:"id"`

	OwnerID uuid.UUID  `db:"owner_id"`
	Owner   *userModel `db:"owner"`

	GameID *uuid.UUID `db:"game_id"`
	Game   *gameModel `db:"game"`

	Name   string          `db:"name"`
	Schema string          `db:"schema"`
	Fields pq.GenericArray `db:"fields"`
}

func (pdf *pdfModel) domain() *domain.PDF {
	if pdf != nil {
		return &domain.PDF{
			ID:      pdf.ID,
			OwnerID: pdf.OwnerID,
			Owner:   pdf.Owner.domain(),
			GameID:  pdf.GameID,
			Game:    pdf.Game.domain(),
			Name:    pdf.Name,
			Schema:  pdf.Schema,
		}
	}
	return nil
}

func (db *Database) InsertPDF(ctx context.Context, pdf *domain.PDF, pageCount int) error {

	hstoreFields := make([]hstore.Hstore, pageCount)
	for i := range pageCount {
		hstoreFields[i].Map = map[string]sql.NullString{}
	}

	// Insert the PDF.
	err := db.tx.QueryRowxContext(ctx,
		`INSERT INTO pdfs (id, owner_id, game_id, name, schema, fields)
			VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		uuid.New(), pdf.OwnerID, pdf.GameID, pdf.Name, pdf.Schema, pq.Array(hstoreFields),
	).Scan(&pdf.ID)

	return errors.Wrap(err, "cannot insert PDF")
}

func (db *Database) GetPDFsForOwner(ctx context.Context, ownerID uuid.UUID, view domain.PDFView) ([]*domain.PDF, error) {

	var joins string

	switch view {
	case domain.PDFViewAll_OwnerInfo:
		joins = `LEFT JOIN users ON users.id = pdfs.owner_id`
	case domain.PDFViewAll_GameInfo:
		joins = `LEFT JOIN games ON games.id = pdfs.game_id`
	}

	// Build a query to select PDFs with the owner-ID.
	query := fmt.Sprintf(
		`SELECT %s FROM pdfs %s WHERE pdfs.owner_id = $1`,
		pdfViewColumns[view], joins,
	)

	// Execute the query.
	var models []*pdfModel
	if err := sqlx.SelectContext(ctx, db.tx, &models, query, ownerID); err != nil {
		return nil, errors.Wrap(err, "cannot select PDFs")
	}

	// Convert each model to a domain.PDF.
	ret := make([]*domain.PDF, len(models))
	for i, m := range models {
		ret[i] = m.domain()
	}

	return ret, nil
}

func (db *Database) GetPDFsForGame(ctx context.Context, gameID uuid.UUID, view domain.PDFView) ([]*domain.PDF, error) {

	var joins string

	switch view {
	case domain.PDFViewAll_OwnerInfo:
		joins = `LEFT JOIN users ON users.id = pdfs.owner_id`
	case domain.PDFViewAll_GameInfo:
		joins = `LEFT JOIN games ON games.id = pdfs.game_id`
	}

	// Build a query to select PDFs with the owner-ID.
	query := fmt.Sprintf(
		`SELECT %s FROM pdfs %s WHERE pdfs.game_id = $1`,
		pdfViewColumns[view], joins,
	)

	// Execute the query.
	var models []*pdfModel
	if err := sqlx.SelectContext(ctx, db.tx, &models, query, gameID); err != nil {
		return nil, errors.Wrap(err, "cannot select PDFs")
	}

	// Convert each model to a domain.PDF.
	ret := make([]*domain.PDF, len(models))
	for i, m := range models {
		ret[i] = m.domain()
	}

	return ret, nil
}

func (db *Database) GetPDF(ctx context.Context, pdfID uuid.UUID, view domain.PDFView) (*domain.PDF, error) {

	var joins string

	switch view {
	case domain.PDFViewAll_OwnerInfo:
		joins = `LEFT JOIN users ON users.id = pdfs.owner_id`
	case domain.PDFViewAll_GameInfo:
		joins = `LEFT JOIN games ON games.id = pdfs.game_id`
	case domain.PDFViewAll_OwnerInfo_GameInfo:
		joins = `LEFT JOIN users ON users.id = pdfs.owner_id LEFT JOIN games ON games.id = pdfs.game_id`
	}

	// Build a query to select a PDF with the PDF-ID.
	query := fmt.Sprintf(
		`SELECT %s FROM pdfs %s WHERE pdfs.id = $1`,
		pdfViewColumns[view], joins,
	)

	// Execute the query.
	var model pdfModel
	if err := sqlx.GetContext(ctx, db.tx, &model, query, pdfID); err != nil {
		if err == sql.ErrNoRows {
			return nil, &domain.NormalError{
				Type:   domain.NETypePDFNotFound,
				Detail: fmt.Sprintf("Cannot find a PDF with the PDF-ID"),
			}
		}

		return nil, errors.Wrap(err, "cannot select PDF")
	}

	return model.domain(), nil
}

func (db *Database) GetPDFFields(ctx context.Context, pdfID uuid.UUID, pageIdx int) (map[string]string, error) {

	// Execute a query to select the fields of a PDF with the PDF-ID.
	var hstorePage hstore.Hstore
	err := db.tx.QueryRowxContext(ctx,
		`SELECT fields[$1] FROM pdfs WHERE id = $2`,
		pageIdx+1, pdfID,
	).Scan(&hstorePage)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &domain.NormalError{
				Type:   domain.NETypePDFNotFound,
				Detail: "Cannot find a PDF with the PDF-ID",
			}
		}

		return nil, errors.Wrap(err, "cannot get PDF fields")
	}

	// Convert the hstore to a map without null-strings.
	ret := make(map[string]string, len(hstorePage.Map))
	for k, v := range hstorePage.Map {
		ret[k] = v.String
	}

	return ret, nil
}

func (db *Database) UpdatePDFField(ctx context.Context, pdfID uuid.UUID, pageIdx int, fieldName, fieldValue string) error {

	// Execute a query to update the fields of a PDF with the PDF-ID.
	_, err := db.tx.ExecContext(ctx,
		`UPDATE pdfs SET fields[$1] = fields[$1] || hstore($2, $3) WHERE id = $4`,
		pageIdx+1, fieldName, fieldValue, pdfID,
	)

	return errors.Wrap(err, "cannot update PDF fields")
}

func (db *Database) DeletePDF(ctx context.Context, pdfID, ownerID uuid.UUID) error {

	// Delete the PDF with the PDF and owner IDs.
	_, err := db.tx.ExecContext(ctx,
		`DELETE FROM pdfs WHERE id = $1 AND owner_id = $2`,
		pdfID, ownerID,
	)

	return errors.Wrap(err, "cannot delete PDF")
}

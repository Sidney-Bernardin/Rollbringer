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
	domain.PDFViewMain:      `pdfs.id, pdfs.owner_id, pdfs.game_id, pdfs.name, pdfs.schema`,
	domain.PDFViewBasicInfo: `pdfs.id, pdfs.owner_id, pdfs.game_id, pdfs.name, pdfs.schema, games.id AS "game.id", COALESCE(games.name, '') AS "game.name"`,
}

type pdfModel struct {
	ID uuid.UUID `db:"id"`

	OwnerID uuid.UUID  `db:"owner_id"`
	Owner   *userModel `db:"owner"`

	GameID *uuid.UUID `db:"game_id"`
	Game   *gameModel `db:"game"`

	Name   string          `db:"name"`
	Schema string          `db:"schema"`
	Pages  pq.GenericArray `db:"pages"`
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

	hstorePages := make([]hstore.Hstore, pageCount)
	for i := range pageCount {
		hstorePages[i].Map = map[string]sql.NullString{}
	}

	// Insert the PDF.
	err := db.tx.QueryRowxContext(ctx,
		`INSERT INTO pdfs (id, owner_id, game_id, name, schema, pages)
			VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		uuid.New(), pdf.OwnerID, pdf.GameID, pdf.Name, pdf.Schema, pq.Array(hstorePages),
	).Scan(&pdf.ID)

	return errors.Wrap(err, "cannot insert pdf")
}

func (db *Database) GetPDFsByOwner(ctx context.Context, ownerID uuid.UUID, view domain.PDFView) ([]*domain.PDF, error) {

	var joins string
	if view == domain.PDFViewBasicInfo {
		joins = `LEFT JOIN games ON pdfs.game_id = games.id`
	}

	// Build a query to select PDFs with the owner-ID.
	query := fmt.Sprintf(
		`SELECT %s FROM pdfs %s WHERE pdfs.owner_id = $1`,
		pdfViewColumns[view], joins,
	)

	// Execute the query.
	var models []*pdfModel
	if err := sqlx.SelectContext(ctx, db.tx, &models, query, ownerID); err != nil {
		return nil, errors.Wrap(err, "cannot select pdfs")
	}

	// Convert each model to domain.PDF.
	ret := make([]*domain.PDF, len(models))
	for i, m := range models {
		ret[i] = m.domain()
	}

	return ret, nil
}

func (db *Database) GetPDF(ctx context.Context, pdfID uuid.UUID, pdfFields, ownerFields, gameFields []string) (*domain.PDF, error) {
	return &domain.PDF{
		ID:      pdfID,
		OwnerID: [16]byte{},
		Owner:   &domain.User{},
		GameID:  nil,
		Game:    &domain.Game{},
		Name:    "",
		Schema:  "",
		Pages:   []map[string]string{},
	}, nil
}

func (db *Database) GetPDFFields(ctx context.Context, pdfID uuid.UUID, pageIdx int) (map[string]string, error) {
	return map[string]string{}, nil
}

func (db *Database) UpdatePDFField(ctx context.Context, pdfID uuid.UUID, pageIdx int, fieldName, fieldValue string) error {
	return nil
}

func (db *Database) DeletePDF(ctx context.Context, pdfID, ownerID uuid.UUID) error {
	return nil
}

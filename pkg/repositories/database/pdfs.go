package database

import (
	"context"
	"fmt"
	"rollbringer/pkg/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type gormPDF struct {
	ID uuid.UUID `gorm:"column:id;default:gen_random_uuid()"`

	OwnerID uuid.UUID `gorm:"column:owner_id"`
	Owner   *gormUser

	GameID uuid.UUID `gorm:"column:game_id"`
	Game   *gormGame

	Name   string                                 `gorm:"column:name"`
	Schema string                                 `gorm:"column:schema"`
	Fields datatypes.JSONSlice[map[string]string] `gorm:"column:fields"`
}

func (pdf *gormPDF) TableName() string { return "pdfs" }

func (pdf *gormPDF) domain() *domain.PDF {
	if pdf != nil {
		return &domain.PDF{
			ID:      pdf.ID,
			OwnerID: pdf.OwnerID,
			Owner:   pdf.Owner.domain(),
			GameID:  pdf.GameID,
			Game:    pdf.Game.domain(),
			Name:    pdf.Name,
			Schema:  pdf.Schema,
			Fields:  pdf.Fields,
		}
	}
	return nil
}

func (db *Database) InsertPDF(ctx context.Context, pdf *domain.PDF, pages int) error {

	pdfModel := gormPDF{
		OwnerID: pdf.OwnerID,
		GameID:  pdf.GameID,
		Name:    pdf.Name,
		Schema:  pdf.Schema,
		Fields:  datatypes.JSONSlice[map[string]string](pdf.Fields),
	}

	err := db.gormDB.WithContext(ctx).Create(&pdfModel).Error
	if err != nil {
		return errors.Wrap(err, "cannot create pdf")
	}

	*pdf = *pdfModel.domain()
	return nil
}

func (db *Database) GetPDFsByOwner(ctx context.Context, ownerID uuid.UUID, pdfFields, ownerFields, gameFields []string) ([]*domain.PDF, error) {

	q := db.gormDB.WithContext(ctx).
		Select(columnFields("pdfs", pdfFields)).
		Where("owner_id = ?", ownerID)

	if ownerFields != nil {
		q = q.Joins("Owner", db.gormDB.Select(ownerFields))
	}

	if gameFields != nil {
		q = q.Joins("Game", db.gormDB.Select(gameFields))
	}

	var pdfModels []gormPDF
	if err := q.Find(&pdfModels).Error; err != nil {
		return nil, errors.Wrap(err, "cannot get pdfs by owner")
	}

	pdfs := make([]*domain.PDF, len(pdfModels))
	for i, m := range pdfModels {
		pdfs[i] = m.domain()
	}

	return pdfs, nil
}

func (db *Database) GetPDFsByGame(ctx context.Context, gameID uuid.UUID, pdfFields, ownerFields, gameFields []string) ([]*domain.PDF, error) {

	q := db.gormDB.WithContext(ctx).
		Select(columnFields("pdfs", pdfFields)).
		Where("game_id = ?", gameID)

	if ownerFields != nil {
		q = q.Joins("Owner", db.gormDB.Select(ownerFields))
	}

	if gameFields != nil {
		q = q.Joins("Game", db.gormDB.Select(gameFields))
	}

	var pdfModels []gormPDF
	if err := q.Find(&pdfModels).Error; err != nil {
		return nil, errors.Wrap(err, "cannot get pdfs by game")
	}

	ret := make([]*domain.PDF, len(pdfModels))
	for i, m := range pdfModels {
		ret[i] = m.domain()
	}

	return ret, nil
}

func (db *Database) GetPDF(ctx context.Context, pdfID uuid.UUID, pdfFields, ownerFields, gameFields []string) (*domain.PDF, error) {

	q := db.gormDB.WithContext(ctx).
		Select(columnFields("pdfs", pdfFields)).
		Where("pdfs.id = ?", pdfID)

	if ownerFields != nil {
		q.Joins("Owner", db.gormDB.Select(ownerFields))
	}

	if gameFields != nil {
		q.Joins("Game", db.gormDB.Select(gameFields))
	}

	var pdfModel gormPDF
	if err := q.First(&pdfModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &domain.ProblemDetail{
				Type:   domain.PDTypePDFNotFound,
				Detail: "Cannot find a pdf with the pdf-ID.",
			}
		}

		return nil, errors.Wrap(err, "cannot get pdf")
	}

	return pdfModel.domain(), nil
}

func (db *Database) GetPDFFields(ctx context.Context, pdfID uuid.UUID, pageIdx int) (map[string]string, error) {

	var fields datatypes.JSONType[map[string]string]
	err := db.gormDB.WithContext(ctx).Model(&gormPDF{}).
		Select(fmt.Sprintf("fields->%v", pageIdx)).
		Where("id = ?", pdfID).
		First(&fields).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &domain.ProblemDetail{
				Type:   domain.PDTypePDFNotFound,
				Detail: "Cannot find a pdf with the pdf-ID.",
			}
		}

		return nil, errors.Wrap(err, "cannot get pdf fields")
	}

	return fields.Data(), nil
}

func (db *Database) UpdatePDFField(ctx context.Context, pdfID uuid.UUID, pageIdx int, fieldName, fieldValue string) error {

	result := db.gormDB.WithContext(ctx).Model(&gormPDF{}).
		Where("id = ?", pdfID).
		Update("fields", datatypes.JSONSet("fields").Set(
			fmt.Sprintf("{%v, %s}", pageIdx, fieldName),
			fieldValue,
		))

	if err := result.Error; err != nil {
		return errors.Wrap(err, "cannot update pdf page field")
	}

	if result.RowsAffected == 0 {
		return &domain.ProblemDetail{
			Type:   domain.PDTypePDFNotFound,
			Detail: "Cannot find a pdf with the pdf-ID.",
		}
	}

	return nil
}

func (db *Database) DeletePDF(ctx context.Context, pdfID, ownerID uuid.UUID) error {
	err := db.gormDB.WithContext(ctx).
		Delete(&gormPDF{
			ID:      pdfID,
			OwnerID: ownerID,
		}).Error

	return errors.Wrap(err, "cannot delete pdf")
}

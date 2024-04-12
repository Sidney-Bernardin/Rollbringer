package service

import (
	"context"

	"rollbringer/pkg/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var PDFSchemaPageNames = map[string][]string{
	"DND_CHARACTER_SHEET": {"main", "info", "spells"},
}

func (svc *Service) CreatePDF(ctx context.Context, session *domain.Session, pdf *domain.PDF) error {

	pages := len(PDFSchemaPageNames[pdf.Schema])

	for range pages {
		pdf.Fields = append(pdf.Fields, map[string]string{})
	}

	err := svc.DB.InsertPDF(ctx, pdf, pages)
	return errors.Wrap(err, "cannot insert pdf")
}

func (svc *Service) GetPDFsByOwner(ctx context.Context, ownerID uuid.UUID, pdfFields, ownerFields, gameFields []string) ([]*domain.PDF, error) {
	pdfs, err := svc.DB.GetPDFsByOwner(ctx, ownerID, pdfFields, ownerFields, gameFields)
	return pdfs, errors.Wrap(err, "cannot get pdfs by owner")
}

func (svc *Service) GetPDFsByGame(ctx context.Context, gameID uuid.UUID, pdfFields, ownerFields, gameFields []string) ([]*domain.PDF, error) {
	pdfs, err := svc.DB.GetPDFsByGame(ctx, gameID, pdfFields, ownerFields, gameFields)
	return pdfs, errors.Wrap(err, "cannot get pdfs by game")
}

func (svc *Service) GetPDF(ctx context.Context, pdfID uuid.UUID, pdfFields, ownerFields, gameFields []string) (*domain.PDF, error) {
	pdf, err := svc.DB.GetPDF(ctx, pdfID, pdfFields, ownerFields, gameFields)
	return pdf, errors.Wrap(err, "cannot get pdf")
}

func (svc *Service) GetPDFFields(ctx context.Context, pdfID uuid.UUID, pageNum int) (map[string]string, error) {
	if pageNum < 1 {
		return nil, &domain.ProblemDetail{
			Type:   domain.PDTypeInvalidPDFPageNumber,
			Detail: "Page number must be greater than zero.",
		}
	}

	fields, err := svc.DB.GetPDFFields(ctx, pdfID, pageNum-1)
	return fields, errors.Wrap(err, "cannot get pdf fields")
}

func (svc *Service) DeletePDF(ctx context.Context, session *domain.Session, pdfID uuid.UUID) error {
	err := svc.DB.DeletePDF(ctx, pdfID, session.UserID)
	return errors.Wrap(err, "cannot delete pdf")
}

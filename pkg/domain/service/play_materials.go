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

func (svc *Service) CreatePDF(ctx context.Context, session *domain.Session, pdf *domain.PDF, view domain.PDFView) error {

	err := svc.DB.InsertPDF(ctx, pdf, len(PDFSchemaPageNames[pdf.Schema]))
	if err != nil {
		return errors.Wrap(err, "cannot insert PDF")
	}

	newPDF, err := svc.DB.GetPDF(ctx, pdf.ID, view)
	if err != nil {
		return errors.Wrap(err, "cannot get PDF")
	}

	*pdf = *newPDF
	return nil
}

func (svc *Service) GetPDF(ctx context.Context, pdfID uuid.UUID, view domain.PDFView) (*domain.PDF, error) {
	pdf, err := svc.DB.GetPDF(ctx, pdfID, view)
	return pdf, errors.Wrap(err, "cannot get pdF")
}

func (svc *Service) GetPDFPage(ctx context.Context, pdfID uuid.UUID, pageNum int) (map[string]string, error) {
	if pageNum < 1 {
		return nil, &domain.ProblemDetail{
			Type:   domain.PDTypeInvalidPDFPageNumber,
			Detail: "Page number must be greater than zero.",
		}
	}

	fields, err := svc.DB.GetPDFPage(ctx, pdfID, pageNum-1)
	return fields, errors.Wrap(err, "cannot get PDF fields")
}

func (svc *Service) DeletePDF(ctx context.Context, session *domain.Session, pdfID uuid.UUID) error {
	err := svc.DB.DeletePDF(ctx, pdfID, session.UserID)
	return errors.Wrap(err, "cannot delete PDF")
}

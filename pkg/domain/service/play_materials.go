package service

import (
	"context"

	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
)

func (svc *Service) CreatePDF(ctx context.Context, userID, schema string) (string, string, error) {
	pdfID, title, err := svc.DB.InsertPDF(ctx, userID, schema)
	return pdfID, title, errors.Wrap(err, "cannot insert PDF")
}

func (svc *Service) GetPDF(ctx context.Context, pdfID string) (*domain.PDF, error) {
	pdf, err := svc.DB.GetPDF(ctx, pdfID)
	return pdf, errors.Wrap(err, "cannot get pdf")
}

func (svc *Service) GetPDFs(ctx context.Context, ownerID string) ([]*domain.PDF, error) {
	pdf, err := svc.DB.GetPDFs(ctx, ownerID)
	return pdf, errors.Wrap(err, "cannot get pdfs")
}

func (svc *Service) DeletePDF(ctx context.Context, pdfID, userID string) error {
	err := svc.DB.DeletePDF(ctx, pdfID, userID)
	return errors.Wrap(err, "cannot delete pdf")
}

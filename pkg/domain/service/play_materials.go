package service

import (
	"context"

	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/repositories/database"
)

var PDFSchemaToPageCount = map[string]int{
	"DND_CHARACTER_SHEET": 3,
	"DND_LEVELING_GUIDE":  1,
}

func (svc *Service) CreatePDF(ctx context.Context, pdf *domain.PDF) (*domain.PDF, error) {

	if len(pdf.Name) < 1 || 30 < len(pdf.Name) {
		return nil, domain.ErrInvalidPDFName
	}

	pdf.Pages = make([]map[string]string, PDFSchemaToPageCount[pdf.Schema])

	err := svc.DB.Transaction(ctx, func(db *database.Database) error {

		// Insert the PDF.
		if err := db.InsertPDF(ctx, pdf); err != nil {
			return errors.Wrap(err, "cannot insert PDF")
		}

		if pdf.GameID == "" {
			return nil
		}

		// Append the PDF to the game.
		if err := db.AppendGamePDF(ctx, pdf.GameID, pdf.ID); err != nil {
			if err == domain.ErrGameNotFound {
				return errors.New("cannot add pdf to game because the game was not found")
			}

			return errors.Wrap(err, "cannot add pdf to game")
		}

		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "transaction failed")
	}

	return pdf, nil
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
	err := svc.DB.Transaction(ctx, func(db *database.Database) error {

		pdf, err := db.GetPDF(ctx, pdfID)
		if err != nil {
			return errors.Wrap(err, "cannot get pdf")
		}

		if pdf.GameID != "" {
			err := db.RemoveGamePDF(ctx, pdf.GameID, pdfID)
			if err != nil && err != domain.ErrGameNotFound {
				return errors.Wrap(err, "cannot remove pdf from game")
			}
		}

		err = db.DeletePDF(ctx, pdfID, userID)
		return errors.Wrap(err, "cannot delete pdf")
	})

	return errors.Wrap(err, "transaction failed")
}

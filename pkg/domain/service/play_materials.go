package service

import (
	"context"

	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
)

func (svc *Service) PlayMaterials(
	ctx context.Context,
	gameID string,
	incomingChan chan *domain.GameEvent,
	outgoingChan chan *domain.GameEvent,
) {
	defer close(outgoingChan)

	// Get the game.
	game, err := svc.GetGame(ctx, gameID)
	if err != nil && errors.Cause(err) != domain.ErrGameNotFound {
		svc.logger.Error().Stack().Err(err).Msg("Cannot get game")
		return
	}

	if game != nil {
		go svc.ps.SubToGame(ctx, game.ID, outgoingChan)
	}

	for {
		select {
		case <-ctx.Done():
			return

		case event := <-incomingChan:
			switch event.Type {

			// case "UPDATE_PDF_PAGE":

			// 	err = svc.db.UpdatePDFPage(ctx, event.PDFID, event.SenderID, event.PageNum, event.PDFFields)
			// 	if err != nil {
			// 		if errors.Cause(err) == domain.ErrPlayMaterialNotFound {
			// 			continue
			// 		}

			// 		svc.logger.Error().Stack().Err(err).Msg("Cannot update PDF")
			// 		return
			// 	}

			// 	if err := svc.ps.PubToGame(ctx, gameID, event); err != nil {
			// 		svc.logger.Error().Stack().Err(err).Msg("Cannot publish game event")
			// 		return
			// 	}
			}
		}
	}
}

func (svc *Service) CreatePDF(ctx context.Context, userID, schema string) (string, string, error) {

	// Insert a new PDF.
	pdfID, title, err := svc.db.InsertPDF(ctx, userID, schema)
	return pdfID, title, errors.Wrap(err, "cannot insert PDF")
}

func (svc *Service) GetPDF(ctx context.Context, ownerID, pdfID string) (*domain.PDF, error) {

	// Get the PDF.
	pdf, err := svc.db.GetPDF(ctx, ownerID, pdfID)
	return pdf, errors.Wrap(err, "cannot get pdf")
}

func (svc *Service) GetPDFs(ctx context.Context, ownerID string) ([]*domain.PDF, error) {

	// Get the PDF.
	pdf, err := svc.db.GetPDFs(ctx, ownerID)
	return pdf, errors.Wrap(err, "cannot get pdfs")
}

func (svc *Service) DeletePDF(ctx context.Context, pdfID, userID string) error {

	// Delete the PDF.
	err := svc.db.DeletePDF(ctx, pdfID, userID)
	return errors.Wrap(err, "cannot delete pdf")
}

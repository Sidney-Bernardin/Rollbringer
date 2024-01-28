package service

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
)

func (svc *Service) PlayMaterials(ctx context.Context, gameID string, incomingChan, outgoingChan chan domain.GameEvent) {

	// Get the game.
	game, err := svc.GetGame(ctx, gameID)
	if err != nil && errors.Cause(err) != domain.ErrGameNotFound {
		svc.logger.Error().Stack().Err(err).Msg("Cannot get game")
		return
	}

	if game != nil {
		go svc.ps.SubToGame(ctx, game.ID, incomingChan)
	}

	for {
		select {
		case <-ctx.Done():
			return

		case event := <-incomingChan:
			switch event["type"] {
			case "PDF_UPDATE":

				content, err := json.Marshal(event)
				if err != nil {
					svc.logger.Error().Stack().Err(err).Msg("Cannot encode game event")
					return
				}

				err = svc.db.UpdatePDF(ctx, event["pdf_id"].(string), event["sender_id"].(string), content)
				if err != nil {
					svc.logger.Error().Stack().Err(err).Msg("Cannot get PDF")
					return
				}
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

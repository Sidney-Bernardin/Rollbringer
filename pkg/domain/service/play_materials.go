package service

import (
	"context"
	"fmt"
	"rollbringer/pkg/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
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
			fmt.Println(event)

			switch event["type"] {
			case "pdf_update":
			}
		}
	}
}

func (svc *Service) CreatePDF(ctx context.Context, userID, schema string) (string, string, error) {

	// Parse the user-ID.
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return "", "", domain.ErrUserNotFound
	}

	// Insert a new PDF.
	id, title, err := svc.db.InsertPDF(ctx, parsedUserID, schema)
	return id.String(), title, errors.Wrap(err, "cannot insert PDF")
}

func (svc *Service) GetPDFs(ctx context.Context, ownerID string) ([]*domain.PDF, error) {

	// Parse the owner-ID.
	parsedOwnerID, err := uuid.Parse(ownerID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// Get the PDF.
	pdf, err := svc.db.GetPDFs(ctx, parsedOwnerID)
	return pdf, errors.Wrap(err, "cannot get pdfs")
}

func (svc *Service) DeletePDF(ctx context.Context, pdfID, userID string) error {

	// Parse the game-ID.
	parsedPDFID, err := uuid.Parse(pdfID)
	if err != nil {
		return domain.ErrPlayMaterialNotFound
	}

	// Parse the user-ID.
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	// Delete the pdf.
	err = svc.db.DeletePDF(ctx, parsedPDFID, parsedUserID)
	return errors.Wrap(err, "cannot delete pdf")
}

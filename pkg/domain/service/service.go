package service

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/repositories/database"
	"rollbringer/pkg/repositories/pubsub"
)

type Service struct {
	DB *database.Database
	PS *pubsub.PubSub

	Logger *zerolog.Logger
}

// DoEvents processes events. outgoingChan closes before returning.
func (svc *Service) DoEvents(
	ctx context.Context,
	gameID string,
	incomingChan chan *domain.Event,
	outgoingChan chan *domain.Event,
) {
	defer close(outgoingChan)

	// Get the game.
	game, err := svc.GetGame(ctx, gameID)
	if err != nil && errors.Cause(err) != domain.ErrGameNotFound {
		svc.Logger.Error().Stack().Err(err).Msg("Cannot get game")
		return
	}

	if game != nil {
		// Subscribe to the game's events.
		go svc.PS.SubToGame(ctx, game.ID, outgoingChan)
	}

	// Process incoming events.
	for {
		select {
		case <-ctx.Done():
			return

		case event := <-incomingChan:
			switch event.Type {

			case "INIT_PDF_PAGE":

				// Get the PDF.
				pdf, err := svc.DB.GetPDF(ctx, event.PDFID)
				if err != nil {
					if errors.Cause(err) == domain.ErrPlayMaterialNotFound {
						continue
					}

					svc.Logger.Error().Stack().Err(err).Msg("Cannot get PDF")
					return
				}

				// Decode the PDF page.
				err = json.Unmarshal([]byte(pdf.Pages[event.PageNum-1]), &event.PDFFields)
				if err != nil {
					svc.Logger.Error().Stack().Err(err).Msg("Cannot unmarshal PDF page")
					return
				}

				outgoingChan <- event

			case "UPDATE_PDF_PAGE":

				// Update the PDF page.
				err = svc.DB.UpdatePDFPage(ctx, event.PDFID, event.PageNum, event.PDFFields)
				if err != nil {
					if errors.Cause(err) == domain.ErrPlayMaterialNotFound {
						continue
					}

					svc.Logger.Error().Stack().Err(err).Msg("Cannot update PDF")
					return
				}

				// Publish the event for the game.
				if err := svc.PS.PubToGame(ctx, gameID, event); err != nil {
					svc.Logger.Error().Stack().Err(err).Msg("Cannot publish event")
					return
				}
			}
		}
	}
}

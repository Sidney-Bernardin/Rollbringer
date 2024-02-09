package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/repositories/database"
	"rollbringer/pkg/repositories/pubsub"
)

type Service struct {
	db *database.Database
	ps *pubsub.PubSub

	logger *zerolog.Logger
}

func New(logger *zerolog.Logger, db *database.Database, ps *pubsub.PubSub) *Service {
	return &Service{
		db:     db,
		ps:     ps,
		logger: logger,
	}
}

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

			case "UPDATE_PDF_PAGE":

				err = svc.db.UpdatePDFPage(ctx, event.PDFID, event.PageNum, event.PDFFields)
				if err != nil {
					if errors.Cause(err) == domain.ErrPlayMaterialNotFound {
						continue
					}

					svc.logger.Error().Stack().Err(err).Msg("Cannot update PDF")
					return
				}

				if err := svc.ps.PubToGame(ctx, gameID, event); err != nil {
					svc.logger.Error().Stack().Err(err).Msg("Cannot publish event")
					return
				}
			}
		}
	}
}

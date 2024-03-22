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
	DB *database.Database
	PS *pubsub.PubSub

	Logger *zerolog.Logger
}

func (svc *Service) GetPlayPage(ctx context.Context, session *domain.Session, gameID string) (page *domain.PlayPage, err error) {

	page = &domain.PlayPage{
		LoggedIn: false,
	}

	// Get the game.
	page.Game, err = svc.GetGame(ctx, gameID)
	if err != nil && !domain.IsProblemDetail(err, domain.PDTypeGameNotFound) {
		return nil, errors.Wrap(err, "cannot get game")
	}

	if session != nil {
		page.LoggedIn = true

		// Get the user.
		page.User, err = svc.DB.GetUser(ctx, session.UserID)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get user")
		}
	}

	return page, nil
}

// DoEvents processes events. errChan closes before returning.
func (svc *Service) DoEvents(ctx context.Context, gameID string, errChan chan error, incomingChan, outgoingChan chan domain.Event) {
	defer close(errChan)

	// Get the game.
	game, err := svc.GetGame(ctx, gameID)
	if err != nil && !domain.IsProblemDetail(err, domain.PDTypeGameNotFound) {
		svc.Logger.Error().Stack().Err(err).Msg("Cannot get game")
		return
	}

	if game != nil {
		// Subscribe to the game's topic.
		go svc.PS.Sub(ctx, game.ID, outgoingChan)
	}

	// Process incoming events.
	for {
		select {
		case <-ctx.Done():
			return

		case e := <-incomingChan:
			switch event := e.(type) {
			case *domain.EventUpdatePDFField:

				if event.PageNum-1 < 0 {
					return
				}

				// Update the PDF page.
				err = svc.DB.UpdatePDFField(ctx, event.PDFID, event.PageNum-1, event.Headers.HXTrigger, event.FieldValue)
				if err != nil {
					errChan <- errors.Wrap(err, "cannot update pdf field")
					continue
				}

				// Publish the event to the game's topic.
				if err := svc.PS.Pub(ctx, gameID, e); err != nil {
					svc.Logger.Error().Stack().Err(err).Msg("Cannot publish event")
					return
				}
			}
		}
	}
}

package service

import (
	"context"

	"github.com/mitchellh/mapstructure"
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
	if err != nil && errors.Cause(err) != domain.ErrGameNotFound {
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

// DoEvents processes events. outgoingChan closes before returning.
func (svc *Service) DoEvents(ctx context.Context, gameID string, incomingChan, outgoingChan chan domain.Event) {
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

		case e := <-incomingChan:
			switch e["OPERATION"] {

			case "UPDATE_PDF_FIELD":

				var event domain.EventUpdatePDFField
				decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
					Result:           &event,
					WeaklyTypedInput: true,
				})

				if err != nil {
					svc.Logger.Error().Stack().Err(err).Msg("Cannot create decoder")
					return
				}

				if err := decoder.Decode(e); err != nil {
					svc.Logger.Error().Stack().Err(err).Msg("Cannot decode event")
					return
				}

				if event.PageNum-1 < 0 {
					return
				}

				// Update the PDF page.
				err = svc.DB.UpdatePDFField(ctx,
					event.PDFID,
					event.PageNum-1,
					event.Headers["HX-Trigger-Name"],
					e[event.Headers["HX-Trigger-Name"]].(string))

				if err != nil {
					if errors.Cause(err) == domain.ErrPlayMaterialNotFound {
						continue
					}

					svc.Logger.Error().Stack().Err(err).Msg("Cannot update PDF")
					return
				}

				// Publish the event.
				if err := svc.PS.PubToGame(ctx, gameID, e); err != nil {
					svc.Logger.Error().Stack().Err(err).Msg("Cannot publish event")
					return
				}
			}
		}
	}
}

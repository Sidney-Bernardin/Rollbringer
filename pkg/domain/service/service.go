package service

import (
	"context"
	"strings"

	"github.com/google/uuid"
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

func (svc *Service) GetPlayPage(ctx context.Context, session *domain.Session, gameID uuid.UUID) (page *domain.PlayPage, err error) {
	page = &domain.PlayPage{}

	if gameID != uuid.Nil {
		page.Game, err = svc.DB.GetGame(ctx, gameID, []string{"id", "title"}, nil)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get game")
		}

		page.GamePDFs, err = svc.DB.GetPDFsByGame(ctx, page.Game.ID, []string{"id", "name", "schema"}, []string{"username"}, nil)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get pdfs by game")
		}
	}

	if session == nil {
		return page, nil
	}

	page.User, err = svc.DB.GetUser(ctx, session.UserID, []string{"id", "username"})
	if err != nil {
		return nil, errors.Wrap(err, "cannot get user")
	}
	page.LoggedIn = true

	page.UserGames, err = svc.DB.GetGamesByHost(ctx, page.User.ID, []string{"id", "title"}, nil)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get games by host")
	}

	page.UserPDFs, err = svc.DB.GetPDFsByOwner(ctx, page.User.ID, []string{"id", "name", "schema"}, nil, []string{"title"})
	if err != nil {
		return nil, errors.Wrap(err, "cannot get pdfs by owner")
	}

	return page, nil
}

// DoEvents processes events. errChan closes before returning.
func (svc *Service) DoEvents(ctx context.Context, gameID uuid.UUID, incomingChan, outgoingChan chan domain.Event, errChan chan error) {
	defer close(errChan)

	var (
		gameCtx       context.Context
		cancelGameCtx context.CancelFunc

		pdfID        uuid.UUID
		pdfCtx       context.Context
		cancelPDFCtx context.CancelFunc = func() {}
	)

	if gameID != uuid.Nil {

		_, err := svc.DB.GetGame(ctx, gameID, []string{"id"}, nil)
		if err != nil {
			errChan <- errors.Wrap(err, "cannot get game")
			return
		}

		gameCtx, cancelGameCtx = context.WithCancel(context.Background())
		defer cancelGameCtx()

		go svc.PS.Sub(gameCtx, gameID.String(), outgoingChan, errChan)
	}

	// Process incoming events.
	for {
		select {
		case <-ctx.Done():
			return

		case e := <-incomingChan:
			switch event := e.(type) {

			case *domain.EventSubToPDFPage:
				cancelPDFCtx()

				pdf, err := svc.DB.GetPDF(ctx, event.PDFID, []string{"id"}, nil, nil)
				if err != nil {
					errChan <- errors.Wrap(err, "cannot get pdf")
					continue
				}
				pdfID = pdf.ID

				pdfCtx, cancelPDFCtx = context.WithCancel(context.Background())
				defer cancelPDFCtx()

				go svc.PS.Sub(pdfCtx, pdfID.String(), outgoingChan, errChan)

			case *domain.EventUpdatePDFField:

				if pdfID == uuid.Nil {
					errChan <- &domain.ProblemDetail{
						Type:   domain.PDTypeNotSubscribedToPDF,
						Detail: "You must be subscribed to a PDF before updating it's field.",
					}
					continue
				}

				if event.PageNum-1 < 0 {
					return
				}

				if strings.Contains(event.FieldName, " ") {
					errChan <- &domain.ProblemDetail{
						Type:   domain.PDTypeInvalidPDFFieldName,
						Detail: "Field name cannot contain spaces.",
					}
					return
				}

				err := svc.DB.UpdatePDFField(ctx, pdfID, event.PageNum-1, event.FieldName, event.FieldValue)
				if err != nil {
					errChan <- errors.Wrap(err, "cannot update pdf field")
					continue
				}

				event.PDFID = pdfID

				if err := svc.PS.Pub(ctx, pdfID.String(), event); err != nil {
					svc.Logger.Error().Stack().Err(err).Msg("Cannot publish event")
					return
				}
			}
		}
	}
}

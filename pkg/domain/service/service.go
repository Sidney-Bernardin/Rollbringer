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

	if session != nil {
		page.User, err = svc.DB.GetUser(ctx, session.UserID, domain.UserViewMain)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get user")
		}

		page.LoggedIn = true

		page.User.PDFs, err = svc.DB.GetPDFsByOwner(ctx, session.UserID, domain.PDFViewBasicInfo)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get user pdfs")
		}

		page.User.HostedGames, err = svc.DB.GetGamesByHost(ctx, session.UserID, domain.GameViewMain)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get user hosted games")
		}

		// page.User.JoinedGames, err = svc.DB.GetGamesJoinedByUser(ctx, session.UserID, domain.GameViewAll, domain.UserViewAll)
		// if err != nil {
		// 	return nil, errors.Wrap(err, "cannot get user joined games")
		// }
	} else {
		page.User = &domain.User{
			Username: "Guest123",
		}
	}

	if gameID != uuid.Nil {
		page.Game, err = svc.DB.GetGame(ctx, gameID, domain.GameViewBasicInfo)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get game")
		}

		if page.Game.HostID == page.User.ID {
			page.IsHost = true
		}

		// page.Game.PDFs, err = svc.DB.GetPDFsByGame(ctx, gameID, domain.GameViewAll)
		// if err != nil {
		// 	return nil, errors.Wrap(err, "cannot get game pdfs")
		// }

		// page.Game.Players, err = svc.DB.GetUsersByJoinedGame(ctx, session.UserID, domain.UserViewAll)
		// if err != nil {
		// 	return nil, errors.Wrap(err, "cannot get game players")
		// }
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

		// TODO: Replace with a GameExists function.
		_, err := svc.DB.GetGame(ctx, gameID, domain.GameViewMain)
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

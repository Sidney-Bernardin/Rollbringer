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

	page.User, err = svc.DB.GetUser(ctx, session.UserID, domain.UserViewAll)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get user")
	}

	page.LoggedIn = true

	page.User.PDFs, err = svc.DB.GetPDFsForOwner(ctx, session.UserID, domain.PDFViewAll_GameInfo)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get user pdfs")
	}

	page.User.HostedGames, err = svc.DB.GetGamesForHost(ctx, session.UserID, domain.GameViewAll)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get user hosted games")
	}

	page.User.JoinedGames, err = svc.DB.GetJoinedGamesForUser(ctx, session.UserID, domain.GameViewAll_HostInfo)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get games joined by user")
	}

	if gameID != uuid.Nil {
		page.Game, err = svc.DB.GetGame(ctx, gameID, domain.GameViewAll)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get game")
		}

		if page.Game.HostID == page.User.ID {
			page.IsHost = true
		}

		page.Game.PDFs, err = svc.DB.GetPDFsForGame(ctx, gameID, domain.PDFViewAll_OwnerInfo)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get game pdfs")
		}

		page.Game.Players, err = svc.DB.GetJoinedUsersForGame(ctx, session.UserID, domain.UserViewAll)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get game players")
		}
	}

	return page, nil
}

// DoEvents processes events. errChan closes before returning.
func (svc *Service) DoEvents(
	ctx context.Context,
	session *domain.Session,
	gameID uuid.UUID,
	incomingChan, outgoingChan chan domain.Event,
	errChan chan error,
) {
	defer close(errChan)

	var (
		gameSubscriptionCtx       context.Context
		cancelGameSubscriptionCtx context.CancelFunc

		pdfID                    uuid.UUID
		pdfSubscriptionCtx       context.Context
		cancelPDFSubscriptionCtx context.CancelFunc = func() {}
	)

	if gameID != uuid.Nil {

		// TODO: Replace with a GameExists function.
		_, err := svc.DB.GetGame(ctx, gameID, domain.GameViewAll)
		if err != nil {
			errChan <- errors.Wrap(err, "cannot get game")
			return
		}

		gameSubscriptionCtx, cancelGameSubscriptionCtx = context.WithCancel(context.Background())
		defer cancelGameSubscriptionCtx()

		go svc.PS.Sub(gameSubscriptionCtx, gameID.String(), outgoingChan, errChan)
	}

	// Process incoming events.
	for {
		select {
		case <-ctx.Done():
			return

		case e := <-incomingChan:

			if e.GetHeaders()["CSRF-Token"] != session.CSRFToken {
				errChan <- &domain.ProblemDetail{
					Type: domain.PDTypeUnauthorized,
				}
				continue
			}

			switch event := e.(type) {

			case *domain.EventSubToPDFPage:
				cancelPDFSubscriptionCtx()

				pdf, err := svc.DB.GetPDF(ctx, event.PDFID, domain.PDFViewAll)
				if err != nil {
					errChan <- errors.Wrap(err, "cannot get pdf")
					continue
				}
				pdfID = pdf.ID

				pdfSubscriptionCtx, cancelPDFSubscriptionCtx = context.WithCancel(context.Background())
				defer cancelPDFSubscriptionCtx()

				go svc.PS.Sub(pdfSubscriptionCtx, pdfID.String(), outgoingChan, errChan)

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

				if event.FieldName == "" || strings.Contains(event.FieldName, " ") {
					errChan <- &domain.ProblemDetail{
						Type:   domain.PDTypeInvalidPDFFieldName,
						Detail: "Field name cannot be empty or contain spaces.",
					}
					return
				}

				err := svc.DB.UpdatePDFPage(ctx, pdfID, event.PageNum-1, event.FieldName, event.FieldValue)
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

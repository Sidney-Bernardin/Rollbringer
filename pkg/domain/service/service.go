package service

import (
	"context"
	"strings"

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
	domain.ParseUUIDs(&gameID)

	page = &domain.PlayPage{
		LoggedIn: false,
	}

	// Get the game.
	page.Game, err = svc.GetGame(ctx, gameID)
	if err != nil && !domain.IsProblemDetail(err, domain.PDTypeGameNotFound) {
		return nil, errors.Wrap(err, "cannot get game")
	}

	if session != nil {
		domain.ParseUUIDs(&session.UserID)
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
func (svc *Service) DoEvents(ctx context.Context, gameID string, incomingChan, outgoingChan chan domain.Event, errChan chan error) {
	defer close(errChan)

	var (
		gameCtx       context.Context
		cancelGameCtx context.CancelFunc

		pdfID        string
		pdfCtx       context.Context
		cancelPDFCtx context.CancelFunc = func() {}
	)

	if gameID != "" {
		domain.ParseUUIDs(&gameID)

		// Get the game.
		_, err := svc.GetGame(ctx, gameID)
		if err != nil {
			errChan <- errors.Wrap(err, "cannot get game")
			return
		}

		gameCtx, cancelGameCtx = context.WithCancel(context.Background())
		defer cancelGameCtx()

		// Subscribe to the game's topic.
		go svc.PS.Sub(gameCtx, gameID, outgoingChan, errChan)
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
				domain.ParseUUIDs(&event.PDFID)

				// Get the PDF.
				pdf, err := svc.DB.GetPDF(ctx, event.PDFID)
				if err != nil {
					errChan <- errors.Wrap(err, "cannot get pdf")
					continue
				}
				pdfID = pdf.ID

				pdfCtx, cancelPDFCtx = context.WithCancel(context.Background())
				defer cancelPDFCtx()

				// Subscribe to the PDF page's topic.
				go svc.PS.Sub(pdfCtx, pdfID, outgoingChan, errChan)

			case *domain.EventUpdatePDFField:

				if pdfID == "" {
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

				// Update the PDF page.
				err := svc.DB.UpdatePDFField(ctx, pdfID, event.PageNum-1, event.FieldName, event.FieldValue)
				if err != nil {
					errChan <- errors.Wrap(err, "cannot update pdf field")
					continue
				}

				event.PDFID = pdfID

				// Publish the event to the PDF's topic.
				if err := svc.PS.Pub(ctx, pdfID, event); err != nil {
					svc.Logger.Error().Stack().Err(err).Msg("Cannot publish event")
					return
				}
			}
		}
	}
}

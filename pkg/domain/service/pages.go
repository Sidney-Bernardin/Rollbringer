package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/repositories/database"
	"rollbringer/pkg/repositories/oauth"
	"rollbringer/pkg/repositories/pubsub"
)

type Service struct {
	DB *database.Database
	PS *pubsub.PubSub
	OA *oauth.OAuth

	Logger *zerolog.Logger

	random *rand.Rand
}

func NewService(
	db *database.Database,
	ps *pubsub.PubSub,
	oa *oauth.OAuth,
	logger *zerolog.Logger,
) *Service {
	return &Service{
		DB:     db,
		PS:     ps,
		OA:     oa,
		Logger: logger,
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
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

		page.Game.Rolls, err = svc.DB.GetRollsForGame(ctx, gameID)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get game rolls")
		}
	}

	return page, nil
}

// DoEvents processes events. errChan closes before returning.
func (svc *Service) DoEvents(
	ctx context.Context,
	session *domain.Session,
	gameID uuid.UUID,
	incomingChan <-chan domain.Event,
	outgoingChan chan domain.Event,
	errChan chan error,
) {
	defer close(errChan)

	var (
		pdfCtx       context.Context
		cancelPDFCtx context.CancelFunc = func() {}
		pdfID        uuid.UUID
	)

	if gameID != uuid.Nil {

		// TODO: Replace with a GameExists function.
		_, err := svc.DB.GetGame(ctx, gameID, domain.GameViewAll)
		if err != nil {
			errChan <- errors.Wrap(err, "cannot get game")
			return
		}

		go svc.PS.Sub(ctx, gameID.String(), outgoingChan, errChan)
	}

	// Process incoming events.
	for {
		select {
		case <-ctx.Done():
			cancelPDFCtx()
			return

		case e := <-incomingChan:
			if err := e.Validate(ctx); err != nil {
				errChan <- errors.Wrap(err, "invalid event")
				continue
			}

			switch event := e.(type) {
			case *domain.EventSubToPDF:
				cancelPDFCtx()
				eventCtx := context.WithValue(ctx, domain.CtxKeyInstance, domain.OperationSubToPDF)

				fields, err := svc.DB.GetPDFFields(eventCtx, event.PDFID, event.PageNum-1)
				if err != nil {
					errChan <- errors.Wrap(err, "cannot get pdf")
					continue
				}

				pdfID = event.PDFID
				pdfCtx, cancelPDFCtx = context.WithCancel(eventCtx)

				topic := fmt.Sprintf("%s_%v", pdfID, event.PageNum)
				go svc.PS.Sub(pdfCtx, topic, outgoingChan, errChan)

				outgoingChan <- &domain.EventPDFFields{
					BaseEvent: domain.BaseEvent{Operation: domain.OperationPDFFields},
					PDFID:     pdfID,
					PageNum:   event.PageNum,
					Fields:    fields,
				}

			case *domain.EventUpdatePDFField:
				eventCtx := context.WithValue(ctx, domain.CtxKeyInstance, domain.OperationUpdatePDFField)

				if pdfID == uuid.Nil {
					errChan <- &domain.NormalError{
						Instance: eventCtx.Value(domain.CtxKeyInstance).(string),
						Type:     domain.NETypeNotSubscribedToPDF,
						Detail:   "You must be subscribed to a PDF before updating it's field.",
					}
					continue
				}

				err := svc.DB.UpdatePDFField(eventCtx, pdfID, event.PageNum-1, event.FieldName, event.FieldValue)
				if err != nil {
					errChan <- errors.Wrap(err, "cannot update pdf field")
					continue
				}

				err = svc.PS.Pub(ctx, pdfID.String(), &domain.EventPDFFields{
					BaseEvent: domain.BaseEvent{Operation: domain.OperationPDFFields},
					PDFID:     pdfID,
					PageNum:   event.PageNum,
					Fields: map[string]string{
						event.FieldName: event.FieldValue,
					},
				})

				if err != nil {
					svc.Logger.Error().Stack().Err(err).Msg("Cannot publish event")
					cancelPDFCtx()
					return
				}

			case *domain.EventCreateRoll:
				eventCtx := context.WithValue(ctx, domain.CtxKeyInstance, domain.OperationCreateRoll)

				roll, err := domain.NewRoll(ctx, svc.random, session.UserID, gameID, event.Dice)
				if err != nil {
					errChan <- errors.Wrap(err, "cannot insert roll")
					continue
				}

				if err := svc.DB.InsertRoll(eventCtx, roll); err != nil {
					errChan <- errors.Wrap(err, "cannot insert roll")
					continue
				}

				err = svc.PS.Pub(ctx, gameID.String(), &domain.EventRoll{
					BaseEvent: domain.BaseEvent{Operation: domain.OperationRoll},
					Roll:      roll,
				})

				if err != nil {
					svc.Logger.Error().Stack().Err(err).Msg("Cannot publish event")
					cancelPDFCtx()
					return
				}
			}
		}
	}
}

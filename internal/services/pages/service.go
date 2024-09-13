package pages

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"rollbringer/internal"
	"rollbringer/internal/config"
	"rollbringer/internal/services"
)

type Service interface {
	services.BaseServicer

	GetSession(ctx context.Context, sessionID uuid.UUID, view internal.SessionView) (*internal.Session, error)
	PlayPage(ctx context.Context, session *internal.Session, gameID uuid.UUID) (*PlayPage, error)
}

type service struct {
	*services.BaseService
}

func NewService(cfg *config.Config, logger *slog.Logger, ps internal.PubSub) Service {
	return &service{
		BaseService: &services.BaseService{
			Config: cfg,
			Logger: logger,
			PS:     ps,
		},
	}
}

func (svc *service) Shutdown() error {
	svc.PS.Close()
	return nil
}

func (svc *service) GetSession(ctx context.Context, sessionID uuid.UUID, sessionView internal.SessionView) (*internal.Session, error) {

	var session internal.Session
	err := svc.PS.Request(ctx, "sessions", &session, &internal.EventWrapper[any]{
		Event: internal.EventGetSessionRequest,
		Payload: internal.GetSessionRequest{
			SessionID:   sessionID,
			SessionView: sessionView,
		},
	})

	if err != nil {
		return nil, errors.Wrap(err, "cannot get session")
	}

	return &session, nil
}

func (svc *service) PlayPage(ctx context.Context, session *internal.Session, gameID uuid.UUID) (*PlayPage, error) {

	var (
		page = &PlayPage{}
		errs = &errgroup.Group{}
	)

	errs.Go(func() error {
		err := svc.PS.Request(ctx, "users", page.User, &internal.EventWrapper[any]{
			Event: internal.EventGetUserRequest,
			Payload: &internal.GetUserRequest{
				UserID: session.UserID,
				View:   "user-all,games-all,pdfs-all",
			},
		})
		return errors.Wrap(err, "cannot get user")
	})

	errs.Go(func() error {
		err := svc.PS.Request(ctx, "games", page.Game, &internal.EventWrapper[any]{
			Event: internal.EventGetGameRequest,
			Payload: &internal.GetGameRequest{
				GameID: gameID,
				View:   "game-all",
			},
		})
		return errors.Wrap(err, "cannot get game")
	})

	if err := errs.Wait(); err != nil {
		return nil, err
	}

	page.IsHost = page.Game.HostID == page.User.ID
	return page, nil
}

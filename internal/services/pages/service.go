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

type PlayPage struct {
	IsHost bool

	User *internal.User
	Game *internal.Game
}

type Service interface {
	services.Servicer

	GetSession(ctx context.Context, sessionID uuid.UUID) (*internal.Session, error)
	PlayPage(ctx context.Context, session *internal.Session, gameID uuid.UUID) (*PlayPage, error)
}

type service struct {
	*services.Service

	ps internal.PubSub
}

func NewService(cfg *config.Config, logger *slog.Logger, ps internal.PubSub) Service {
	return &service{
		Service: &services.Service{
			Config: cfg,
			Logger: logger,
		},
		ps: ps,
	}
}

func (svc *service) Shutdown() error {
	svc.ps.Close()
	return nil
}

func (svc *service) GetSession(ctx context.Context, sessionID uuid.UUID) (*internal.Session, error) {
	res, err := svc.ps.Request(ctx, "sessions", &internal.EventGetSession{
		BaseEvent: internal.BaseEvent{Type: internal.ETGetSession},
		SessionID: sessionID,
	})

	if err != nil {
		return nil, errors.Wrap(err, "cannot authenticate")
	}

	event, _ := res.(*internal.EventSession)
	return &event.Session, nil
}

func (svc *service) PlayPage(ctx context.Context, session *internal.Session, gameID uuid.UUID) (*PlayPage, error) {

	var (
		page  = &PlayPage{}
		group = &errgroup.Group{}
	)

	group.Go(func() error {
		res, err := svc.ps.Request(ctx, "users", &internal.EventGetUser{
			UserID: session.UserID,
			View:   internal.UserViewAll,
		})

		if err != nil {
			return errors.Wrap(err, "cannot get user")
		}

		event, _ := res.(*internal.EventUser)
		page.User = &event.User
		return nil
	})

	group.Go(func() error {
		res, err := svc.ps.Request(ctx, "games", &internal.EventGetGame{
			GameID: gameID,
			View:   internal.GameViewAll,
		})

		if err != nil {
			return errors.Wrap(err, "cannot get game")
		}

		event, _ := res.(*internal.EventGame)
		page.Game = &event.Game
		return nil
	})

	if err := group.Wait(); err != nil {
		return nil, err
	}

	page.IsHost = page.Game.HostID == page.User.ID
	return page, nil
}

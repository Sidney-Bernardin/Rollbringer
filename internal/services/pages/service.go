package pages

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"rollbringer/internal"
	"rollbringer/internal/config"
	"rollbringer/internal/repositories/pubsub"
	"rollbringer/internal/services"
)

type PlayPage struct {
	IsHost bool

	User *internal.User
	Game *internal.Game
}

type Service interface {
	services.Servicer

	PlayPage(ctx context.Context, session *internal.Session, gameID uuid.UUID) (*PlayPage, error)
	Authenticate(ctx context.Context, sessionID uuid.UUID, csrfToken string) (*internal.Session, error)
}

type service struct {
	*services.Service

	ps *pubsub.PubSub
}

func NewService(cfg *config.Config, logger *slog.Logger, ps *pubsub.PubSub) Service {
	return &service{
		Service: &services.Service{
			Config: cfg,
			Logger: logger.With("component", "pages_service"),
		},
		ps: ps,
	}
}

func (svc *service) Authenticate(ctx context.Context, sessionID uuid.UUID, csrfToken string) (*internal.Session, error) {
	return nil, nil
}

func (svc *service) PlayPage(ctx context.Context, session *internal.Session, gameID uuid.UUID) (*PlayPage, error) {

	var (
		page    = &PlayPage{}
		errChan = make(chan error)
	)

	go func() {
		res, err := svc.ps.Request(ctx, fmt.Sprintf("users.%s", session.UserID), nil)
		errChan <- errors.Wrap(err, "cannot get user")

		err = json.Unmarshal(res, &page.User)
		errChan <- errors.Wrap(err, "cannot JSON decode user")
	}()

	go func() {
		res, err := svc.ps.Request(ctx, fmt.Sprintf("games.%s", session.UserID), nil)
		errChan <- errors.Wrap(err, "cannot get game")

		err = json.Unmarshal(res, &page.Game)
		errChan <- errors.Wrap(err, "cannot JSON decode game")
	}()

	for range 4 {
		if err := <-errChan; err != nil {
			return nil, err
		}
	}

	page.IsHost = page.Game.HostID == page.User.ID
	return page, nil
}

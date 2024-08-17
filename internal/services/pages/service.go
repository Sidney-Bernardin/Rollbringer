package service

import (
	"context"
	"encoding"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"rollbringer/internal"
	"rollbringer/internal/config"
	"rollbringer/internal/repositories/pubsub"
)

type PlayPage struct {
	IsHost bool

	User *internal.User
	Game *internal.Game
}

type PagesService interface {
	PlayPage(ctx context.Context, session *internal.Session, gameID uuid.UUID) (*PlayPage, error)
	Authenticate(ctx context.Context, sessionID uuid.UUID, csrfToken string) (*internal.Session, error)
}

type service struct {
	cfg    *config.Config
	logger *slog.Logger

	ps *pubsub.PubSub
}

func New(cfg *config.Config, logger *slog.Logger, ps *pubsub.PubSub) PagesService {
	return &service{
		cfg:    cfg,
		logger: logger,
		ps:     ps,
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

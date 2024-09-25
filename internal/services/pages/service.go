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
			PubSub: ps,
		},
	}
}

func (svc *service) Shutdown() error {
	svc.PubSub.Close()
	return nil
}

func (svc *service) PlayPage(ctx context.Context, session *internal.Session, gameID uuid.UUID) (*PlayPage, error) {

	var (
		page          = &PlayPage{}
		errs, errsCtx = errgroup.WithContext(ctx)
	)

	errs.Go(func() error {
		err := svc.PubSub.Request(errsCtx, "users", &page.User, &internal.EventWrapper[any]{
			Event: internal.EventGetUserRequest,
			Payload: &internal.GetUserRequest{
				UserID:    session.UserID,
				ViewQuery: "user-all,games-all,pdfs-all",
			},
		})
		return errors.Wrap(err, "cannot get user")
	})

	if gameID != uuid.Nil {
		errs.Go(func() error {
			err := svc.PubSub.Request(errsCtx, "games", &page.Game, &internal.EventWrapper[any]{
				Event: internal.EventGetGameRequest,
				Payload: &internal.GetGameRequest{
					GameID:    gameID,
					ViewQuery: "game-all,users-all,pdfs-all,rolls-all",
				},
			})

			if err != nil {
				return errors.Wrap(err, "cannot get game")
			}

			page.IsHost = page.Game.HostID == session.UserID
			return nil
		})
	}

	if err := errs.Wait(); err != nil {
		return nil, err
	}

	return page, nil
}

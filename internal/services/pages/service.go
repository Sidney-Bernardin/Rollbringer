package pages

import (
	"context"
	"log/slog"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"rollbringer/internal"
	"rollbringer/internal/config"
	"rollbringer/internal/services"
)

type Service interface {
	services.BaseServicer

	PlayPage(ctx context.Context, page *internal.PlayPage) error
}

type service struct {
	*services.BaseService
}

func NewService(cfg *config.Config, logger *slog.Logger, ps internal.PubSub) Service {
	svc := &service{
		BaseService: &services.BaseService{
			Config: cfg,
			Logger: logger,
			PubSub: ps,
		},
	}

	return svc
}

func (svc *service) Shutdown() error {
	svc.PubSub.Close()
	return nil
}

func (svc *service) PlayPage(ctx context.Context, page *internal.PlayPage) error {
	errs, errsCtx := errgroup.WithContext(ctx)

	if page.Game != nil {
		errs.Go(func() error {
			var res []*internal.User
			err := svc.PubSub.Request(errsCtx, "users.users", &res, &internal.EventWrapper[any]{
				Event: internal.EventGetUsersForGameRequest,
				Payload: internal.GetUsersForGameRequest{
					GameID: page.Game.ID,
				},
			})

			if err != nil {
				return errors.Wrap(err, "cannot get users for game")
			}

			page.Game.Users = res
			return nil
		})
	}

	errs.Go(func() error {
		var res internal.PlayPage
		err := svc.PubSub.Request(errsCtx, "games.pages", &res, &internal.EventWrapper[any]{
			Event:   internal.EventPlayPage,
			Payload: page,
		})

		if err != nil {
			return errors.Wrap(err, "cannot get game data for play-page")
		}

		page.Session.User = res.Session.User
		if page.Game != nil {
			page.Game.Name = res.Game.Name
			page.Game.HostID = res.Game.HostID
			page.Game.PDFs = res.Game.PDFs
			page.Game.Rolls = res.Game.Rolls
			page.Game.ChatMessages = res.Game.ChatMessages
		}

		return nil
	})

	return errs.Wait()
}

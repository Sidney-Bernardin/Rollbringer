package games

import (
	"context"

	"github.com/pkg/errors"

	"rollbringer/internal"
)

func (svc *service) doSubscriptions(ctx context.Context) {
	var errChan = make(chan error)

	go svc.ps.Subscribe(ctx, "games", errChan, func(e internal.Event, subject []string) (internal.Event, *internal.ProblemDetail) {
		switch event := e.(type) {
		case *internal.EventGetGame:
			game, err := svc.getGame(ctx, event.GameID, event.View)
			if err != nil {
				return nil, svc.HandleError(ctx, errors.Wrap(err, "cannot get game"))
			}

			return &internal.EventGame{
				BaseEvent: internal.BaseEvent{Type: internal.ETGame},
				Game:      *game,
			}, nil

		default:
			return nil, nil
		}
	})

	for {
		select {
		case <-ctx.Done():
			return

		case err := <-errChan:
			svc.HandleError(ctx, errors.Wrap(err, "subscription error"))
		}
	}
}

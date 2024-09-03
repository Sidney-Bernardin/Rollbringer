package games

import (
	"context"

	"github.com/pkg/errors"

	"rollbringer/internal"
)

func (svc *service) Run() error {
	errChan := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	select {
	case <-ctx.Done():
		return nil
	case err := <-errChan:
		return err
	}
}

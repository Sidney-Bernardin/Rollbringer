package games

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"rollbringer/internal"
)

func (svc *service) Run() error {
	errChan := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go svc.doGamesSubscription(ctx, errChan)

	select {
	case <-ctx.Done():
		return nil
	case err := <-errChan:
		return err
	}
}

func (svc *service) doGamesSubscription(ctx context.Context, errChan chan<- error) {
	svc.PS.Subscribe(ctx, "games", errChan, func(incoming *internal.EventWrapper[[]byte]) *internal.EventWrapper[any] {

		var payload any
		switch incoming.Event {
		case internal.EventGetGameRequest:
			payload = internal.GetGameRequest{}
		default:
			return nil
		}

		if err := json.Unmarshal(incoming.Payload, &payload); err != nil {
			return &internal.EventWrapper[any]{
				Event:   internal.EventError,
				Payload: internal.HandleError(ctx, svc.Logger, errors.Wrap(err, "cannot JSON decode incoming payload")),
			}
		}

		switch incoming := payload.(type) {
		case internal.GetGameRequest:
			game, err := svc.getGame(ctx, &incoming)
			if err != nil {
				return &internal.EventWrapper[any]{
					Event:   internal.EventError,
					Payload: internal.HandleError(ctx, svc.Logger, errors.Wrap(err, "cannot get game")),
				}
			}

			return &internal.EventWrapper[any]{
				Event:   internal.EventGame,
				Payload: game,
			}

		default:
			return nil
		}
	})
}

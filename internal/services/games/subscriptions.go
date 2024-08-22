package games

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"rollbringer/internal"
)

func subscriptionError()

func (svc *service) doSubscriptions(ctx context.Context) {
	var errChan = make(chan error)

	go svc.ps.Subscribe(ctx, "games.*.*", errChan, func(e internal.Event, subject []string) internal.Event {

		var (
			gameID, _ = uuid.Parse(subject[1])
			view, _   = internal.GameViews[subject[2]]
		)

		game, err := svc.db.GameGet(ctx, gameID, view)
		if err != nil {
			return &internal.EventError{
				BaseEvent:     internal.BaseEvent{Type: internal.ETError},
				ProblemDetail: internal.HandleError(ctx, svc.Logger, errors.Wrap(err, "cannot get game")),
			}
		}

		return &internal.EventGame{
			BaseEvent: internal.BaseEvent{Type: internal.ETGame},
			Game:      game,
		}
	})

	for {
		select {
		case <-ctx.Done():
			return

		case err := <-errChan:
			internal.HandleError(ctx, svc.Logger, errors.Wrap(err, "subscription error"))
		}
	}
}

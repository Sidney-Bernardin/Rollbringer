package users

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"rollbringer/internal"
)

func (svc *service) Listen() error {
	errChan := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go svc.doUsersSubscription(ctx, errChan)

	select {
	case <-ctx.Done():
		return nil
	case err := <-errChan:
		return err
	}
}

func (svc *service) doUsersSubscription(ctx context.Context, errChan chan<- error) {
	svc.PS.Subscribe(ctx, "users", errChan, func(incoming *internal.EventWrapper[[]byte]) *internal.EventWrapper[any] {

		var payload any
		switch incoming.Event {
		case internal.EventGetUserRequest:
			payload = internal.GetUserRequest{}
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
		case internal.GetUserRequest:
			user, err := svc.getUser(ctx, incoming.UserID, incoming.View)
			if err != nil {
				return &internal.EventWrapper[any]{
					Event:   internal.EventError,
					Payload: internal.HandleError(ctx, svc.Logger, errors.Wrap(err, "cannot get game")),
				}
			}

			return &internal.EventWrapper[any]{
				Event:   internal.EventUser,
				Payload: user,
			}

		case *internal.AuthenticateUserRequest:
			session, err := svc.authenticate(ctx, incoming.SessionID, incoming.CSRFToken)
			if err != nil {
				return &internal.EventWrapper[any]{
					Event:   internal.EventError,
					Payload: internal.HandleError(ctx, svc.Logger, errors.Wrap(err, "cannot authenticate")),
				}
			}

			return &internal.EventWrapper[any]{
				Event:     internal.EventSession,
				Payload:   session,
			}

		default:
			return nil
		}
	})
}

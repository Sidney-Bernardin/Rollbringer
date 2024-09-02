package users

import (
	"context"

	"github.com/pkg/errors"

	"rollbringer/internal"
)

func (svc *service) doSubscriptions(ctx context.Context) {
	var errChan = make(chan error)

	go svc.ps.Subscribe(ctx, "users", errChan, func(e internal.Event, subject []string) (internal.Event, *internal.ProblemDetail) {
		switch event := e.(type) {
		case *internal.EventGetUser:
			user, err := svc.getUser(ctx, event.UserID, event.View)
			if err != nil {
				return nil, svc.HandleError(ctx, errors.Wrap(err, "cannot get user"))
			}

			return &internal.EventUser{
				BaseEvent: internal.BaseEvent{Type: internal.ETUser},
				User:      *user,
			}, nil

		case *internal.EventAuthenticate:
			session, err := svc.authenticate(ctx, event.SessionID, event.CSRFToken)
			if err != nil {
				return nil, svc.HandleError(ctx, errors.Wrap(err, "cannot authenticate"))
			}

			return &internal.EventSession{
				BaseEvent: internal.BaseEvent{Type: internal.ETSession},
				Session:   *session,
			}, nil

		default:
			return nil, nil
		}
	})

	go svc.ps.Subscribe(ctx, "sessions", errChan, func(e internal.Event, subject []string) (internal.Event, *internal.ProblemDetail) {
		switch event := e.(type) {
		case *internal.EventGetSession:
			session, err := svc.getSession(ctx, event.SessionID, event.View)
			if err != nil {
				return nil, svc.HandleError(ctx, errors.Wrap(err, "cannot get session"))
			}

			return &internal.EventSession{
				BaseEvent: internal.BaseEvent{Type: internal.ETSession},
				Session:   *session,
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

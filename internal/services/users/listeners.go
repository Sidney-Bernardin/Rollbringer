package users

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"rollbringer/internal"
)

func (svc *service) Listen() error {
	errs, ctx := errgroup.WithContext(context.Background())

	errs.Go(func() error {
		return svc.subToUsers(ctx)
	})

	errs.Go(func() error {
		return svc.subToSessions(ctx)
	})

	return errs.Wait()
}

func (svc *service) subToUsers(ctx context.Context) error {
	err := svc.PubSub.Subscribe(ctx, "users", func(req *internal.EventWrapper[[]byte]) *internal.EventWrapper[any] {

		var payload any
		switch req.Event {
		case internal.EventGetUserRequest:
			payload = &internal.GetUserRequest{}
		case internal.EventAuthenticateUserRequest:
			payload = &internal.AuthenticateUserRequest{}
		default:
			return &internal.EventWrapper[any]{
				Event: internal.EventError,
				Payload: internal.NewProblemDetail(ctx, internal.PDOpts{
					Type:   internal.PDTypeInvalidEvent,
					Detail: "The given evnet is not valid.",
					Extra: map[string]any{
						"event": req.Event,
					},
				}),
			}
		}

		instanceCtx := context.WithValue(ctx, internal.CtxKeyInstance, req.Event)

		if err := json.Unmarshal(req.Payload, &payload); err != nil {
			return &internal.EventWrapper[any]{
				Event: internal.EventError,
				Payload: internal.NewProblemDetail(instanceCtx, internal.PDOpts{
					Type:   internal.PDTypeInvalidJSON,
					Detail: err.Error(),
				}),
			}
		}

		switch payload := payload.(type) {
		case *internal.GetUserRequest:
			user, err := svc.getUser(instanceCtx, payload.UserID, payload.ViewQuery)
			if err != nil {
				return &internal.EventWrapper[any]{
					Event:   internal.EventError,
					Payload: internal.HandleError(instanceCtx, svc.Logger, errors.Wrap(err, "cannot get user")),
				}
			}

			return &internal.EventWrapper[any]{
				Event:   internal.EventUser,
				Payload: user,
			}

		case *internal.AuthenticateUserRequest:
			session, err := svc.authenticate(ctx, payload.SessionID, payload.CSRFToken)
			if err != nil {
				return &internal.EventWrapper[any]{
					Event:   internal.EventError,
					Payload: internal.HandleError(instanceCtx, svc.Logger, errors.Wrap(err, "cannot authenticate")),
				}
			}

			return &internal.EventWrapper[any]{
				Event:   internal.EventSession,
				Payload: session,
			}

		default:
			return &internal.EventWrapper[any]{
				Event:   internal.EventError,
				Payload: internal.HandleError(ctx, svc.Logger, internal.SvrErrUnknownEvent),
			}
		}
	})

	return errors.Wrap(err, "cannot subscribe to users")
}

func (svc *service) subToSessions(ctx context.Context) error {
	err := svc.PubSub.Subscribe(ctx, "sessions", func(req *internal.EventWrapper[[]byte]) *internal.EventWrapper[any] {

		var payload any
		switch req.Event {
		case internal.EventGetSessionRequest:
			payload = &internal.GetSessionRequest{}
		default:
			return &internal.EventWrapper[any]{
				Event: internal.EventError,
				Payload: internal.NewProblemDetail(ctx, internal.PDOpts{
					Type:   internal.PDTypeInvalidEvent,
					Detail: "The given evnet is not valid.",
					Extra: map[string]any{
						"event": req.Event,
					},
				}),
			}
		}

		instanceCtx := context.WithValue(ctx, internal.CtxKeyInstance, req.Event)

		if err := json.Unmarshal(req.Payload, &payload); err != nil {
			return &internal.EventWrapper[any]{
				Event: internal.EventError,
				Payload: internal.NewProblemDetail(instanceCtx, internal.PDOpts{
					Type:   internal.PDTypeInvalidJSON,
					Detail: err.Error(),
				}),
			}
		}

		switch payload := payload.(type) {
		case *internal.GetSessionRequest:
			session, err := svc.getSession(instanceCtx, payload.SessionID, payload.ViewQuery)
			if err != nil {
				return &internal.EventWrapper[any]{
					Event:   internal.EventError,
					Payload: internal.HandleError(instanceCtx, svc.Logger, errors.Wrap(err, "cannot get session")),
				}
			}

			return &internal.EventWrapper[any]{
				Event:   internal.EventUser,
				Payload: session,
			}

		default:
			return &internal.EventWrapper[any]{
				Event:   internal.EventError,
				Payload: internal.HandleError(ctx, svc.Logger, internal.SvrErrUnknownEvent),
			}
		}
	})

	return errors.Wrap(err, "cannot subscribe to sessions")
}

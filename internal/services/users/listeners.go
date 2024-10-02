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

	return errs.Wait()
}

func (svc *service) subToUsers(ctx context.Context) error {
	err := svc.PubSub.Subscribe(ctx, "users.users", func(req *internal.EventWrapper[[]byte]) (*internal.EventWrapper[any], error) {

		var payload any
		switch req.Event {
		case internal.EventGetUsersForGameRequest:
			payload = &internal.GetUsersForGameRequest{}
		case internal.EventAuthenticateUserRequest:
			payload = &internal.AuthenticateRequest{}
		default:
			return nil, internal.NewProblemDetail(ctx, internal.PDOpts{
				Type:   internal.PDTypeInvalidEvent,
				Detail: "The given evnet is not valid.",
				Extra: map[string]any{
					"event": req.Event,
				},
			})
		}

		instanceCtx := context.WithValue(ctx, internal.CtxKeyInstance, req.Event)

		if err := json.Unmarshal(req.Payload, &payload); err != nil {
			return nil, internal.NewProblemDetail(instanceCtx, internal.PDOpts{
				Type:   internal.PDTypeInvalidJSON,
				Detail: err.Error(),
			})
		}

		switch payload := payload.(type) {
		case *internal.GetUsersForGameRequest:
			user, err := svc.getUsersForGame(instanceCtx, payload.GameID)
			return &internal.EventWrapper[any]{
				Event:   internal.EventUser,
				Payload: user,
			}, errors.Wrap(err, "cannot get users for game")

		case *internal.AuthenticateRequest:
			session, err := svc.authenticate(instanceCtx, payload.SessionID, payload.SessionView, payload.CheckCSRFToken, payload.CSRFToken)
			return &internal.EventWrapper[any]{
				Event:   internal.EventSession,
				Payload: session,
			}, errors.Wrap(err, "cannot authenticate")

		default:
			return nil, internal.HandleError(ctx, svc.Logger, errors.New("invalid users.users event"))
		}
	})

	return errors.Wrap(err, "cannot subscribe to users")
}

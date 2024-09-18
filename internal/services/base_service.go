package services

import (
	"context"
	"log/slog"

	"rollbringer/internal"
	"rollbringer/internal/config"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type BaseServicer interface {
	Listen() error
	Shutdown() error

	Authenticate(ctx context.Context, sessionID uuid.UUID, csrfToken string) (*internal.Session, error)
	GetSession(ctx context.Context, sessionID uuid.UUID, viewQuery string) (*internal.Session, error)
}

type BaseService struct {
	Config *config.Config
	Logger *slog.Logger

	PubSub internal.PubSub
}

func (svc *BaseService) Listen() error   { return nil }
func (svc *BaseService) Shutdown() error { return nil }

func (svc *BaseService) Authenticate(ctx context.Context, sessionID uuid.UUID, csrfToken string) (*internal.Session, error) {

	var session internal.Session
	err := svc.PubSub.Request(ctx, "users", &session, &internal.EventWrapper[any]{
		Event: internal.EventAuthenticateUserRequest,
		Payload: internal.AuthenticateUserRequest{
			SessionID: sessionID,
			CSRFToken: csrfToken,
		},
	})

	if err != nil {
		return nil, errors.Wrap(err, "cannot authenticate")
	}

	return &session, nil
}

func (svc *BaseService) GetSession(ctx context.Context, sessionID uuid.UUID, viewQuery string) (*internal.Session, error) {

	var session internal.Session
	err := svc.PubSub.Request(ctx, "sessions", &session, &internal.EventWrapper[any]{
		Event: internal.EventGetSessionRequest,
		Payload: internal.GetSessionRequest{
			SessionID: sessionID,
			ViewQuery: viewQuery,
		},
	})

	if err != nil {
		return nil, errors.Wrap(err, "cannot get session")
	}

	return &session, nil
}

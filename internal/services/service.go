package services

import (
	"context"
	"log/slog"

	"rollbringer/internal"
	"rollbringer/internal/config"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Servicer interface {
	HandleError(ctx context.Context, err error) *internal.ProblemDetail
	Authenticate(ctx context.Context, sessionID uuid.UUID, csrfToken string) (*internal.Session, error)
}

type Service struct {
	Config *config.Config
	Logger *slog.Logger

	ps internal.PubSub
}

func (svc *Service) Authenticate(ctx context.Context, sessionID uuid.UUID, csrfToken string) (*internal.Session, error) {
	res, err := svc.ps.Request(ctx, "users", &internal.EventAuthenticate{
		BaseEvent: internal.BaseEvent{Type: internal.ETGetSession},
		SessionID: sessionID,
		CSRFToken: csrfToken,
	})

	if err != nil {
		return nil, errors.Wrap(err, "cannot authenticate")
	}

	event, _ := res.(*internal.EventSession)
	return &event.Session, nil
}

func (svc *Service) HandleError(ctx context.Context, err error) *internal.ProblemDetail {
	pd, ok := errors.Cause(err).(*internal.ProblemDetail)
	if !ok {
		svc.Logger.Error("Server error", "err", err.Error())
		pd = internal.NewProblemDetail(ctx, internal.PDOpts{
			Type: internal.PDTypeServerError,
		})
	}
	return pd
}

package domain

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

var (
	SlogLevelTrace slog.Level = -8
	SlogLevelFatal slog.Level = 12
)

type IService interface {
	Run(context.Context) error
	Shutdown(context.Context) error

	GetSession(ctx context.Context, sessionID uuid.UUID) (*Session, error)
}

type Service struct {
	Config *Config
	Logger *slog.Logger

	PubSub PubSubRepository
}

func (svc *Service) Run(context.Context) error { return nil }

func (svc *Service) Shutdown(ctx context.Context) error {
	svc.PubSub.Close()
	return nil
}

func (svc *Service) GetSession(ctx context.Context, sessionID uuid.UUID) (*Session, error) {

	var session *Session
	_, err := svc.PubSub.Request(ctx, "accounts", &session, &Event{
		Operation: OperationGetSessionRequest,
		Payload:   GetSessionRequest{sessionID},
	})

	if err != nil {
		return nil, Wrap(err, "cannot get session", nil)
	}

	return session, nil
}

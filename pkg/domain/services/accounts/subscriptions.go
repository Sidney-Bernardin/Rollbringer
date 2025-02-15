package service

import (
	"context"
	"log/slog"
	"rollbringer/pkg/domain"
)

func (svc *accountsService) subAccounts(ctx context.Context, e *domain.Event) *domain.Event {
	switch p := e.Payload.(type) {

	case *domain.GetSessionRequest:
		session, err := svc.GetSession(ctx, p.SessionID)
		if err != nil {
			return &domain.Event{
				Operation: domain.OperationError,
				Payload:   domain.HandleError(ctx, svc.Logger, slog.LevelError, domain.Wrap(err, "cannot get session", nil)),
			}
		} else {
			return &domain.Event{
				Operation: domain.OperationSession,
				Payload:   session,
			}
		}

	default:
		return nil
	}
}

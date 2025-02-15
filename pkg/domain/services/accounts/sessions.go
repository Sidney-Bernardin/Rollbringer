package service

import (
	"context"
	"rollbringer/pkg/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (svc *accountsService) GetSession(ctx context.Context, sessionID uuid.UUID) (*domain.Session, error) {
	session, err := svc.accountsDBRepo.SessionGet(ctx, "id", sessionID)
	if err != nil {
		if errors.Cause(err) == domain.ErrNotFound {
			return nil, domain.UserErr(ctx, domain.UsrErrTypeRecordNotFound, "Cannot find a session with the given session-ID.", nil)
		}

		return nil, domain.Wrap(err, "cannot get session", nil)
	}

	return session, nil
}

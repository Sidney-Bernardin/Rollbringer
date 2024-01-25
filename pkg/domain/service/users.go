package service

import (
	"context"
	"rollbringer/pkg/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (svc *Service) LoginUser(ctx context.Context, googleID string) (string, error) {

	// Login the user.
	sessionID, err := svc.db.Login(ctx, googleID)
	return sessionID.String(), errors.Wrap(err, "cannot login user")
}

func (svc *Service) GetUser(ctx context.Context, userID string) (*domain.User, error) {

	// Parse the user-ID.
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// Get the user.
	user, err := svc.db.GetUser(ctx, parsedUserID)
	return user, errors.Wrap(err, "cannot get user")
}

func (svc *Service) GetSession(ctx context.Context, sessionID string) (*domain.Session, error) {

	// Parse the session-ID.
	parsedSessionID, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	// Get the session.
	session, err := svc.db.GetSession(ctx, parsedSessionID)
	return session, errors.Wrap(err, "cannot get session")
}

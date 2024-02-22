package service

import (
	"context"
	"rollbringer/pkg/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (svc *Service) Login(ctx context.Context, googleID string) (*domain.Session, error) {

	// Create a user.
	user := &domain.User{
		GoogleID: googleID,
		Username: "new-user_123",
	}

	// Insert the user.
	if err := svc.DB.InsertUser(ctx, user); err != nil {
		return nil, errors.Wrap(err, "cannot insert user")
	}

	// Create a session for the user.
	session := &domain.Session{
		CSRFToken: uuid.New().String(),
		UserID:    user.ID,
	}

	// Upsert the session.
	if err := svc.DB.UpsertSession(ctx, session); err != nil {
		return nil, errors.Wrap(err, "cannot upsert session")
	}

	return session, nil
}

func (svc *Service) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	user, err := svc.DB.GetUser(ctx, userID)
	return user, errors.Wrap(err, "cannot get user")
}

func (svc *Service) GetSession(ctx context.Context, sessionID string) (*domain.Session, error) {
	session, err := svc.DB.GetSession(ctx, sessionID)
	return session, errors.Wrap(err, "cannot get session")
}

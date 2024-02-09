package service

import (
	"context"
	"rollbringer/pkg/domain"

	"github.com/pkg/errors"
)

func (svc *Service) LoginUser(ctx context.Context, googleID string) (string, error) {
	sessionID, err := svc.db.Login(ctx, googleID)
	return sessionID, errors.Wrap(err, "cannot login user")
}

func (svc *Service) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	user, err := svc.db.GetUser(ctx, userID)
	return user, errors.Wrap(err, "cannot get user")
}

func (svc *Service) GetSession(ctx context.Context, sessionID string) (*domain.Session, error) {
	session, err := svc.db.GetSession(ctx, sessionID)
	return session, errors.Wrap(err, "cannot get session")
}

package users

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"rollbringer/internal"
	"rollbringer/internal/config"
	"rollbringer/internal/repositories/oauth"
	"rollbringer/internal/services"
)

type Service interface {
	services.Servicer

	StartLogin() (consentURL, state, codeVerifier string)
	FinishLogin(ctx context.Context, stateA, stateB, code, codeVerifier string) (*internal.Session, error)
}

type service struct {
	*services.Service

	ps internal.PubSub
	db internal.UsersDatabase
	oa internal.OAuth
}

func NewService(
	cfg *config.Config,
	logger *slog.Logger,
	ps internal.PubSub,
	db internal.UsersDatabase,
	oa internal.OAuth,
) Service {
	svc := &service{
		Service: &services.Service{
			Config: cfg,
			Logger: logger,
		},
		ps: ps,
		db: db,
		oa: oa,
	}

	return svc
}

func (svc *service) Shutdown() error {
	svc.ps.Close()
	err := svc.db.Close()
	return errors.Wrap(err, "cannot close database")
}

func (svc *service) StartLogin() (consentURL, state, codeVerifier string) {
	state = oauth.NewOauthState()
	codeVerifier = svc.oa.GenerateCodeVerifier()
	return svc.oa.GetConsentURL(state, codeVerifier), state, codeVerifier
}

func (svc *service) FinishLogin(ctx context.Context, stateA, stateB, code, codeVerifier string) (*internal.Session, error) {
	userInfo, err := svc.oa.AuthenticateConsent(ctx, stateA, stateB, code, codeVerifier)
	if err != nil {
		return nil, errors.Wrap(err, "cannot authenticate consent")
	}

	user := &internal.User{
		GoogleID: &userInfo.GoogleID,
		Username: userInfo.GivenName,
	}

	if err := svc.db.UserInsert(ctx, user); err != nil {
		return nil, errors.Wrap(err, "cannot insert user")
	}

	session := &internal.Session{
		UserID:    user.ID,
		CSRFToken: uuid.New().String(),
	}

	if err := svc.db.SessionUpsert(ctx, session); err != nil {
		return nil, errors.Wrap(err, "cannot insert session")
	}

	return session, nil
}

func (svc *service) getUser(ctx context.Context, userID uuid.UUID, view internal.UserView) (*internal.User, error) {
	user, err := svc.db.UserGet(ctx, userID, view)
	return user, errors.Wrap(err, "cannot get user")
}

func (svc *service) getSession(ctx context.Context, sessionID uuid.UUID, view internal.SessionView) (*internal.Session, error) {
	session, err := svc.db.SessionGet(ctx, sessionID, view)
	return session, errors.Wrap(err, "cannot get session")
}

func (svc *service) authenticate(ctx context.Context, sessionID uuid.UUID, csrfToken string) (*internal.Session, error) {
	session, err := svc.db.SessionGet(ctx, sessionID, internal.SessionViewAll)
	if err != nil {
		if internal.IsDetailed(err, internal.PDTypeSessionNotFound) {
			return nil, internal.NewProblemDetail(ctx, internal.PDOpts{
				Type: internal.PDTypeUnauthorized,
			})
		}

		return nil, errors.Wrap(err, "cannot get session")
	}

	if session.CSRFToken != csrfToken {
		return nil, internal.NewProblemDetail(ctx, internal.PDOpts{
			Type: internal.PDTypeUnauthorized,
		})
	}

	return session, nil
}

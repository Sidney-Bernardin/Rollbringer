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
	services.BaseServicer

	GoogleStartLogin() (consentURL, state, codeVerifier string)
	GoogleFinishLogin(ctx context.Context, stateA, stateB, code, codeVerifier string) (*internal.Session, error)
}

type service struct {
	*services.BaseService

	schema internal.UsersSchema
	oauth  internal.OAuth
}

func NewService(
	cfg *config.Config,
	logger *slog.Logger,
	ps internal.PubSub,
	schema internal.UsersSchema,
	oauth internal.OAuth,
) Service {
	svc := &service{
		BaseService: &services.BaseService{
			Config: cfg,
			Logger: logger,
			PubSub: ps,
		},
		schema: schema,
		oauth:  oauth,
	}

	return svc
}

func (svc *service) Shutdown() error {
	svc.PubSub.Close()
	err := svc.schema.Close()
	return errors.Wrap(err, "cannot close database")
}

func (svc *service) GoogleStartLogin() (consentURL, state, codeVerifier string) {
	state = oauth.NewOauthState()
	codeVerifier = svc.oauth.GenerateCodeVerifier()
	return svc.oauth.GetConsentURL(state, codeVerifier), state, codeVerifier
}

func (svc *service) GoogleFinishLogin(ctx context.Context, stateA, stateB, code, codeVerifier string) (*internal.Session, error) {
	userInfo, err := svc.oauth.AuthenticateConsent(ctx, stateA, stateB, code, codeVerifier)
	if err != nil {
		return nil, errors.Wrap(err, "cannot authenticate consent")
	}

	user := &internal.User{
		Username:      userInfo.GivenName,
		GoogleID:      &userInfo.GoogleID,
		GooglePicture: &userInfo.Picture,
	}

	if err := svc.schema.UserInsert(ctx, user); err != nil {
		return nil, errors.Wrap(err, "cannot insert user")
	}

	session := &internal.Session{
		UserID:    user.ID,
		CSRFToken: uuid.New().String(),
	}

	if err := svc.schema.SessionUpsert(ctx, session); err != nil {
		return nil, errors.Wrap(err, "cannot upsert session")
	}

	return session, nil
}

func (svc *service) getUsersForGame(ctx context.Context, gameID uuid.UUID) ([]*internal.User, error) {
	users, err := svc.schema.UsersGetByGame(ctx, gameID)
	return users, errors.Wrap(err, "cannot get users by game")
}

func (svc *service) authenticate(ctx context.Context, sessionID uuid.UUID, sessionView internal.SessionView, checkCSRFToken bool, csrfToken string) (*internal.Session, error) {
	session, err := svc.schema.SessionGet(ctx, sessionID, sessionView)
	if err != nil {
		if internal.IsDetailed(err, internal.PDTypeSessionNotFound) {
			return nil, internal.NewProblemDetail(ctx, internal.PDOpts{
				Type: internal.PDTypeUnauthorized,
			})
		}

		return nil, errors.Wrap(err, "cannot get session")
	}

	if checkCSRFToken && session.CSRFToken != csrfToken {
		return nil, internal.NewProblemDetail(ctx, internal.PDOpts{
			Type: internal.PDTypeUnauthorized,
		})
	}

	return session, nil
}

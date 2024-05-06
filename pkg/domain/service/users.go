package service

import (
	"context"
	"rollbringer/pkg/domain"
	"rollbringer/pkg/repositories/oauth"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (svc *Service) StartLogin() (consentURL, state, codeVerifier string) {
	state = oauth.NewOauthState()
	codeVerifier = svc.OA.NewCodeVerifier()
	return svc.OA.GetConsentURL(state, codeVerifier), state, codeVerifier
}

func (svc *Service) FinishLogin(ctx context.Context, stateA, stateB, code, codeVerifier string) (*domain.Session, error) {

	claims, err := svc.OA.AuthenticateConsent(ctx, stateA, stateB, code, codeVerifier)
	if err != nil {
		return nil, errors.Wrap(err, "cannot authenticate consent")
	}

	user := &domain.User{
		GoogleID: &claims.Subject,
		Username: claims.GivenName,
	}

	if err := svc.DB.InsertUser(ctx, user); err != nil {
		return nil, errors.Wrap(err, "cannot insert user")
	}

	session := &domain.Session{
		UserID:    user.ID,
		CSRFToken: uuid.New().String(),
	}

	if err := svc.DB.UpsertSession(ctx, session); err != nil {
		return nil, errors.Wrap(err, "cannot insert session")
	}

	return session, nil
}

func (svc *Service) Authenticate(ctx context.Context, sessionID uuid.UUID, checkCSRFToken bool, csrfToken string) (*domain.Session, error) {

	session, err := svc.DB.GetSession(ctx, sessionID, domain.SessionViewAll)
	if err != nil {
		if domain.IsProblemDetail(err, domain.PDTypeSessionNotFound) {
			return nil, &domain.ProblemDetail{
				Type: domain.PDTypeUnauthorized,
			}
		}

		return nil, errors.Wrap(err, "cannot get session")
	}

	if checkCSRFToken && session.CSRFToken != csrfToken {
		return nil, &domain.ProblemDetail{
			Type: domain.PDTypeUnauthorized,
		}
	}

	return session, nil
}

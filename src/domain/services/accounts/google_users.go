package accounts

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"rollbringer/src/domain"
)

type GoogleUser struct {
	GoogleID       string
	GivenName      string
	Email          string
	ProfilePicture string
}

func (svc *service) GoogleLogin(ctx context.Context, oauthCode string, newAccount bool) (uuid.UUID, error) {

	googleUser, err := svc.google.GetGoogleUser(ctx, oauthCode)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "google cannot get google-user")
	}

	if !newAccount {

		sessionID, err := svc.database.GoogleSignin(ctx, googleUser)
		if errors.Is(err, domain.ErrNoEntitiesEffected) {
			return uuid.Nil, &domain.ExternalError{Type: ExternalErrorTypeProviderNotLinked, Msg: "The Google account is not linked with a Rollbringer account."}
		}

		return sessionID, errors.Wrap(err, "database cannot signin")
	}

	sessionID, err := svc.database.GoogleSignup(ctx, googleUser, &User{
		ID:             uuid.New(),
		GoogleID:       &googleUser.GoogleID,
		Username:       Username(googleUser.GivenName),
		ProfilePicture: googleUser.ProfilePicture,
	})

	if errors.Is(err, domain.ErrEntityConflict) {
		return uuid.Nil, &domain.ExternalError{Type: ExternalErrorTypeProviderNotLinked, Msg: "The Google account is already linked with a Rollbringer account."}
	}

	return sessionID, errors.Wrap(err, "database cannot signup")
}

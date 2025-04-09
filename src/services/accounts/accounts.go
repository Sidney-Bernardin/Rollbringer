package accounts

import (
	"context"

	"rollbringer/src"
	"rollbringer/src/services/accounts/models"

	"github.com/pkg/errors"
)

const (
	ExternalErrorTypeProviderNotLinked     src.ExternalErrorType = "provider_not_linked"
	ExternalErrorTypeProviderAlreadyLinked src.ExternalErrorType = "provider_already_linked"
)

type Service interface {
	GoogleLogin(ctx context.Context, oauthCode string, newAccount bool) (sessionID *src.UUID, err error)
	SpotifyLogin(ctx context.Context, oauthCode string, newAccount bool) (sessionID *src.UUID, err error)
	Auth(ctx context.Context, sessionID string, csrfToken *string) (*models.Session, error)
}

type service struct {
	config *src.Config

	db      Database
	google  Google
	spotify Spotify
}

func NewService(config *src.Config, db Database, google Google, spotify Spotify) Service {
	return &service{config, db, google, spotify}
}

func (svc *service) Run(ctx context.Context) error {
	return nil
}

func (svc *service) GoogleLogin(ctx context.Context, oauthCode string, newAccount bool) (sessionID *src.UUID, err error) {

	// Get the google-user from Google.
	googleUser, err := svc.google.GetGoogleUser(ctx, oauthCode)
	if err != nil {
		return nil, errors.Wrap(err, "google cannot get google-user")
	}

	if !newAccount {

		// Signin.
		sessionID, err = svc.db.GoogleSignin(ctx, googleUser)
		return sessionID, errors.Wrap(err, "database cannot signin")
	}

	// Create a user.
	user, err := models.NewUser(googleUser, nil)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create user")
	}

	// Signup.
	sessionID, err = svc.db.GoogleSignup(ctx, user)
	return sessionID, errors.Wrap(err, "database cannot signup")
}

func (svc *service) SpotifyLogin(ctx context.Context, oauthCode string, newAccount bool) (sessionID *src.UUID, err error) {

	// Get the spotify-user from Spotify.
	spotifyUser, err := svc.spotify.GetSpotifyUser(ctx, oauthCode)
	if err != nil {
		return nil, errors.Wrap(err, "spotify cannot get spotify-user")
	}

	if !newAccount {

		// Signin.
		sessionID, err = svc.db.SpotifySignin(ctx, spotifyUser)
		return sessionID, errors.Wrap(err, "database cannot signin")
	}

	// Create a user.
	user, err := models.NewUser(nil, spotifyUser)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create user")
	}

	// Signup.
	sessionID, err = svc.db.SpotifySignup(ctx, user)
	return sessionID, errors.Wrap(err, "database cannot signup")
}

func (svc *service) Auth(ctx context.Context, sessionIDStr string, csrfToken *string) (*models.Session, error) {

	// Parse the sessionID.
	sessionID, err := src.ParseUUID(sessionIDStr)
	if err != nil {
		return nil, nil
	}

	// Get the session.
	var session *models.Session
	if csrfToken == nil {
		session, err = svc.db.GetSessionByID(ctx, sessionID)
	} else {
		session, err = svc.db.GetSessionByIDAndCSRFToken(ctx, sessionID, models.CSRFToken(*csrfToken))
	}

	if err != nil {
		if errors.Is(err, src.ErrEntityNotFound) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "cannot get session")
	}

	return session, nil
}

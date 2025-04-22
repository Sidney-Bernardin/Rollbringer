package accounts

import (
	"context"

	"github.com/pkg/errors"

	"rollbringer/src"
	"rollbringer/src/services/accounts/models"
)

func (svc *service) GoogleLogin(ctx context.Context, oauthCode string, newAccount bool) (*src.UUID, error) {

	// Get the google-user from Google.
	googleUser, err := svc.google.GetGoogleUser(ctx, oauthCode)
	if err != nil {
		return nil, errors.Wrap(err, "google cannot get google-user")
	}

	if !newAccount {

		// Signin.
		sessionID, err := svc.database.GoogleSignin(ctx, googleUser)
		if errors.Is(err, src.ErrNoEntitiesEffected) {
			return nil, &src.ExternalError{Type: ExternalErrorTypeProviderNotLinked, Msg: "The Google account is not linked with a Rollbringer account."}
		}

		return sessionID, errors.Wrap(err, "database cannot signin")
	}

	// Create a user.
	user, err := models.NewUser(googleUser, nil)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create user")
	}

	// Signup.
	sessionID, err := svc.database.GoogleSignup(ctx, user)
	if errors.Is(err, src.ErrEntityConflict) {
		return nil, &src.ExternalError{Type: ExternalErrorTypeProviderNotLinked, Msg: "The Google account is already linked with a Rollbringer account."}
	}

	return sessionID, errors.Wrap(err, "database cannot signup")
}

func (svc *service) SpotifyLogin(ctx context.Context, oauthCode string, newAccount bool) (*src.UUID, error) {

	// Get the spotify-user from Spotify.
	spotifyUser, err := svc.spotify.GetSpotifyUser(ctx, oauthCode)
	if err != nil {
		return nil, errors.Wrap(err, "spotify cannot get spotify-user")
	}

	if !newAccount {

		// Signin.
		sessionID, err := svc.database.SpotifySignin(ctx, spotifyUser)
		if errors.Is(err, src.ErrNoEntitiesEffected) {
			return nil, &src.ExternalError{Type: ExternalErrorTypeProviderNotLinked, Msg: "The Spotify account is not linked with a Rollbringer account."}
		}
		return sessionID, errors.Wrap(err, "database cannot signin")
	}

	// Create a user.
	user, err := models.NewUser(nil, spotifyUser)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create user")
	}

	// Signup.
	sessionID, err := svc.database.SpotifySignup(ctx, user)
	if errors.Is(err, src.ErrEntityConflict) {
		return nil, &src.ExternalError{Type: ExternalErrorTypeProviderNotLinked, Msg: "The Google account is already linked with a Rollbringer account."}
	}
	return sessionID, errors.Wrap(err, "database cannot signup")
}

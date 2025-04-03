package accounts

import (
	"context"

	"github.com/pkg/errors"

	"rollbringer/src"
	"rollbringer/src/domain"
)

type User struct {
	UserID domain.UUID

	GoogleID   *string
	GoogleUser *GoogleUser

	SpotifyID   *string
	SpotifyUser *SpotifyUser

	Username       Username
	ProfilePicture string
}

func newUser(googleUser *GoogleUser, spotifyUser *SpotifyUser) (user *User, err error) {
	user = &User{
		UserID:         domain.NewUUID(),
		ProfilePicture: "/static/favicon.png",
	}

	if googleUser != nil {
		user.GoogleID = &googleUser.GoogleID
		user.GoogleUser = googleUser
		user.Username = Username(googleUser.GivenName)
		user.ProfilePicture = googleUser.ProfilePicture
	} else if spotifyUser != nil {
		user.SpotifyID = &spotifyUser.SpotifyID
		user.SpotifyUser = spotifyUser
		user.Username = Username(spotifyUser.DisplayName)
		if spotifyUser.ProfilePicture != nil {
			user.ProfilePicture = *spotifyUser.ProfilePicture
		}
	} else {
		return nil, &src.ExternalError{Type: ExternalErrorTypeUserWithoutProviders}
	}

	return user, nil
}

type GoogleUser struct {
	GoogleID string

	GivenName      string
	Email          string
	ProfilePicture string
}

type SpotifyUser struct {
	SpotifyID string

	DisplayName    string
	Email          string
	ProfilePicture *string
}

/////

type Username string

func ParseUsername(str string) (Username, error) {
	if len(str) == 0 || 25 < len(str) {
		return "", &src.ExternalError{
			Type:        ExternalErrorTypeUsernameInvalid,
			Description: "Must be between 1 and 25 characters",
			Details:     map[string]any{"username": str},
		}
	}

	return Username(str), nil
}

/////

func (svc *service) GoogleLogin(ctx context.Context, oauthCode string, createNewAccount bool) (sessionID domain.UUID, err error) {

	// Get the google-user from Google.
	googleUser, err := svc.google.GetGoogleUser(ctx, oauthCode)
	if err != nil {
		return domain.UUID{}, errors.Wrap(err, "cannot get google-user from Google")
	}

	if !createNewAccount {

		// Signin.
		sessionID, err = svc.db.GoogleSignin(ctx, googleUser)
		return sessionID, errors.Wrap(err, "cannot signin")
	}

	// Create a user.
	user, err := newUser(googleUser, nil)
	if err != nil {
		return domain.UUID{}, errors.Wrap(err, "cannot create user")
	}

	// Signup.
	sessionID, err = svc.db.GoogleSignup(ctx, user)
	return sessionID, errors.Wrap(err, "cannot signup")
}

func (svc *service) SpotifyLogin(ctx context.Context, oauthCode string, createNewAccount bool) (sessionID domain.UUID, err error) {

	// Get the spotify-user from Spotify.
	spotifyUser, err := svc.spotify.GetSpotifyUser(ctx, oauthCode)
	if err != nil {
		return domain.UUID{}, errors.Wrap(err, "cannot get spotify-user from Spotify")
	}

	if !createNewAccount {

		// Signin.
		sessionID, err = svc.db.SpotifySignin(ctx, spotifyUser)
		return sessionID, errors.Wrap(err, "cannot signin")
	}

	// Create a user.
	user, err := newUser(nil, spotifyUser)
	if err != nil {
		return domain.UUID{}, errors.Wrap(err, "cannot create user")
	}

	// Signup.
	sessionID, err = svc.db.SpotifySignup(ctx, user)
	return sessionID, errors.Wrap(err, "cannot signup")
}

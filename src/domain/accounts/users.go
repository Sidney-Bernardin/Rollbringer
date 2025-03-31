package accounts

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"rollbringer/src"
)

type User struct {
	UserID uuid.UUID

	GoogleUser  *GoogleUser
	SpotifyUser *SpotifyUser

	Username Username
}

func newUser(username string, googleUser *GoogleUser, spotifyUser *SpotifyUser) (user *User, err error) {
	if googleUser == nil && spotifyUser == nil {
		return nil, &src.ExternalError{Type: ExternalErrorTypeUserWithoutProviders}
	}

	user = &User{
		UserID:      uuid.New(),
		GoogleUser:  googleUser,
		SpotifyUser: spotifyUser,
	}

	user.Username, err = ParseUsername(username)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse username")
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
			Attrs:       map[string]any{"username": str},
		}
	}

	return Username(str), nil
}

/////

func (svc *service) GoogleLogin(ctx context.Context, oauthState string, createNewAccount bool) (uuid.UUID, error) {
	googleUser, err := svc.google.GetGoogleUser(ctx, oauthState)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "cannot get google-user from Google")
	}

	var sessionID uuid.UUID
	if !createNewAccount {
		sessionID, err = svc.db.GoogleSignin(ctx, googleUser)
		return sessionID, errors.Wrap(err, "cannot signin")
	}

	user, err := newUser(googleUser.GivenName, googleUser, nil)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "cannot create user")
	}

	sessionID, err = svc.db.GoogleSignup(ctx, user)
	return sessionID, errors.Wrap(err, "cannot login")
}

func (svc *service) SpotifyLogin(ctx context.Context, oauthState string, createNewAccount bool) (uuid.UUID, error) {
	spotifyUser, err := svc.spotify.GetSpotifyUser(ctx, oauthState)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "cannot get spotify-user from Google")
	}

	var sessionID uuid.UUID
	if !createNewAccount {
		sessionID, err = svc.db.SpotifySignin(ctx, spotifyUser)
		return sessionID, errors.Wrap(err, "cannot signin")
	}

	user, err := newUser(spotifyUser.DisplayName, nil, spotifyUser)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "cannot create user")
	}

	sessionID, err = svc.db.GoogleSignup(ctx, user)
	return sessionID, errors.Wrap(err, "cannot login")
}

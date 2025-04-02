package accounts

import (
	"context"

	"github.com/google/uuid"
)

type (
	Database interface {
		DatabaseCommands
		DatabaseQueries
	}

	DatabaseCommands interface {
		GoogleSignup(ctx context.Context, user *User) (sessionID uuid.UUID, err error)
		GoogleSignin(ctx context.Context, googleUser *GoogleUser) (sessionID uuid.UUID, err error)

		SpotifySignup(ctx context.Context, user *User) (sessionID uuid.UUID, err error)
		SpotifySignin(ctx context.Context, spotifUser *SpotifyUser) (sessionID uuid.UUID, err error)
	}

	DatabaseQueries interface{}
)

type Google interface {
	ConsentURL() (consentURL string, state string)
	GetGoogleUser(ctx context.Context, oauthCode string) (*GoogleUser, error)
}

type Spotify interface {
	ConsentURL() (consentURL string, state string)
	GetSpotifyUser(ctx context.Context, oauthCode string) (*SpotifyUser, error)
}

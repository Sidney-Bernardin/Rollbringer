package accounts

import (
	"context"
	"rollbringer/src/domain"
)

type (
	Database interface {
		DatabaseCommands
		DatabaseQueries
	}

	DatabaseCommands interface {
		GoogleSignup(ctx context.Context, user *User) (sessionID domain.UUID, err error)
		GoogleSignin(ctx context.Context, googleUser *GoogleUser) (sessionID domain.UUID, err error)

		SpotifySignup(ctx context.Context, user *User) (sessionID domain.UUID, err error)
		SpotifySignin(ctx context.Context, spotifUser *SpotifyUser) (sessionID domain.UUID, err error)
	}

	DatabaseQueries interface {
		SessionGetByID(ctx context.Context, view any, sessionID domain.UUID) error
		SessionGetByIDAndCSRFToken(ctx context.Context, view any, sessionID domain.UUID, csrfToken CSRFToken) error
	}
)

type Google interface {
	ConsentURL() (consentURL string, state string)
	GetGoogleUser(ctx context.Context, oauthCode string) (*GoogleUser, error)
}

type Spotify interface {
	ConsentURL() (consentURL string, state string)
	GetSpotifyUser(ctx context.Context, oauthCode string) (*SpotifyUser, error)
}

package accounts

import (
	"context"
	"rollbringer/src"
	"rollbringer/src/services/accounts/models"
)

type (
	Database interface {
		DatabaseCommands
		DatabaseQueries
	}

	DatabaseCommands interface {
		GoogleSignup(ctx context.Context, user *models.User) (sessionID *src.UUID, err error)
		GoogleSignin(ctx context.Context, googleUser *models.GoogleUser) (sessionID *src.UUID, err error)

		SpotifySignup(ctx context.Context, user *models.User) (sessionID *src.UUID, err error)
		SpotifySignin(ctx context.Context, spotifUser *models.SpotifyUser) (sessionID *src.UUID, err error)
	}

	DatabaseQueries interface {
		GetSessionByID(ctx context.Context, sessionID src.UUID) (*models.Session, error)
		GetSessionByIDAndCSRFToken(ctx context.Context, sessionID src.UUID, csrfToken models.CSRFToken) (*models.Session, error)
	}
)

type Google interface {
	ConsentURL() (consentURL string, state string)
	GetGoogleUser(ctx context.Context, oauthCode string) (*models.GoogleUser, error)
}

type Spotify interface {
	ConsentURL() (consentURL string, state string)
	GetSpotifyUser(ctx context.Context, oauthCode string) (*models.SpotifyUser, error)
}

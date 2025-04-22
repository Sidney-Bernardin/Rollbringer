package accounts

import (
	"context"

	"rollbringer/src"
	"rollbringer/src/services"
	"rollbringer/src/services/accounts/models"
)

const (
	ExternalErrorTypeProviderNotLinked     src.ExternalErrorType = "provider_not_linked"
	ExternalErrorTypeProviderAlreadyLinked src.ExternalErrorType = "provider_already_linked"
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
		GetUsersByRoomID(ctx context.Context, roomID src.UUID) ([]*models.User, error)
		GetUsersByRoomIDs(ctx context.Context, roomIDs ...src.UUID) (map[src.UUID][]*models.User, error)
	}

	Google interface {
		ConsentURL() (consentURL string, state string)
		GetGoogleUser(ctx context.Context, oauthCode string) (*models.GoogleUser, error)
	}

	Spotify interface {
		ConsentURL() (consentURL string, state string)
		GetSpotifyUser(ctx context.Context, oauthCode string) (*models.SpotifyUser, error)
	}
)

type Service interface {
	GoogleLogin(ctx context.Context, oauthCode string, newAccount bool) (sessionID *src.UUID, err error)
	SpotifyLogin(ctx context.Context, oauthCode string, newAccount bool) (sessionID *src.UUID, err error)
	Auth(ctx context.Context, sessionID string, csrfToken *string) (*models.Session, error)
}

type service struct {
	config *src.Config

	broker   services.Broker
	database Database
	google   Google
	spotify  Spotify
}

func NewService(config *src.Config, broker services.Broker, database Database, google Google, spotify Spotify) Service {
	return &service{config, broker, database, google, spotify}
}

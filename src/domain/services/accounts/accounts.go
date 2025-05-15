package accounts

import (
	"context"
	"rollbringer/src"
	"rollbringer/src/domain"

	"github.com/google/uuid"
)

const (
	ExternalErrorTypeProviderNotLinked     domain.ExternalErrorType = "provider-not-linked"
	ExternalErrorTypeProviderAlreadyLinked domain.ExternalErrorType = "provider-already-linked"
)

type Service interface {
	GoogleLogin(ctx context.Context, oauthCode string, newAccount bool) (sessionID uuid.UUID, err error)
	SpotifyLogin(ctx context.Context, oauthCode string, newAccount bool) (sessionID uuid.UUID, err error)
	GetUserByUserID(ctx context.Context, userID uuid.UUID) (*User, error)
}

type (
	Database interface {
		DatabaseCommon

		GoogleSignup(ctx context.Context, googleUser *GoogleUser, user *User) (sessionID uuid.UUID, err error)
		GoogleSignin(ctx context.Context, googleUser *GoogleUser) (sessionID uuid.UUID, err error)

		SpotifySignup(ctx context.Context, spotifUser *SpotifyUser, user *User) (sessionID uuid.UUID, err error)
		SpotifySignin(ctx context.Context, spotifUser *SpotifyUser) (sessionID uuid.UUID, err error)
	}

	DatabaseCommon interface {
		GetSessionBySessionID(ctx context.Context, sessionID uuid.UUID) (*Session, error)
		GetSessionBySessionIDAndCSRFToken(ctx context.Context, sessionID uuid.UUID, csrfToken string) (*Session, error)

		GetUserByUserID(ctx context.Context, userID uuid.UUID) (*User, error)
		GetUsersByUserIDs(ctx context.Context, userIDs ...uuid.UUID) ([]*User, error)
	}
)

type (
	Google interface {
		ConsentURL() (consentURL string, state string)
		GetGoogleUser(ctx context.Context, oauthCode string) (*GoogleUser, error)
	}

	Spotify interface {
		ConsentURL() (consentURL string, state string)
		GetSpotifyUser(ctx context.Context, oauthCode string) (*SpotifyUser, error)
	}
)

type service struct {
	config *src.Config

	publicBroker domain.Broker

	db      Database
	google  Google
	spotify Spotify
}

func NewService(config *src.Config, publicBroker domain.Broker, database Database, google Google, spotify Spotify) Service {
	return &service{config, publicBroker, database, google, spotify}
}

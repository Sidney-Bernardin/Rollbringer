package accounts

import (
	"context"

	"rollbringer/src"

	"github.com/google/uuid"
)

type OAuthProvider int

const (
	ExternalErrorTypeUserWithoutProviders = "user_without_providers"
	ExternalErrorTypeUsernameInvalid      = "username_invalid"
	ExternalErrorTypeUsernameTaken        = "username_taken"
)

type Service interface {
	GoogleLogin(ctx context.Context, oauthState string, createNewAccount bool) (sessionID uuid.UUID, err error)
	SpotifyLogin(ctx context.Context, oauthState string, createNewAccount bool) (sessionID uuid.UUID, err error)
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

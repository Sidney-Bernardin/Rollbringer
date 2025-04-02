package accounts

import (
	"context"

	"rollbringer/src"

	"github.com/google/uuid"
)

type OAuthProvider int

const (
	ExternalErrorTypeUnauthorized src.ExternalErrorType = "unauthorized"

	ExternalErrorTypeUserWithoutProviders src.ExternalErrorType = "user_without_providers"
	ExternalErrorTypeUsernameInvalid      src.ExternalErrorType = "username_invalid"
	ExternalErrorTypeUsernameTaken        src.ExternalErrorType = "username_taken"

	ExternalErrorTypeProviderNotLinked     src.ExternalErrorType = "provider_not_linked"
	ExternalErrorTypeProviderAlreadyLinked src.ExternalErrorType = "provider_already_linked"
)

type Service interface {
	GoogleLogin(ctx context.Context, oauthCode string, createNewAccount bool) (sessionID uuid.UUID, err error)
	SpotifyLogin(ctx context.Context, oauthCode string, createNewAccount bool) (sessionID uuid.UUID, err error)
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

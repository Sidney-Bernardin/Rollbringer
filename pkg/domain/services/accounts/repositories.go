package service

import (
	"context"
	"rollbringer/pkg/domain"

	"golang.org/x/oauth2"
)

type AccountsDatabaseRepository interface {
	Close() error
	Transaction(context.Context, func(AccountsDatabaseRepository) error) error

	UserInsert(ctx context.Context, user *domain.User) error
	UserGet(ctx context.Context, key string, value any) (*domain.User, error)

	GoogleUserInsert(ctx context.Context, googleUser *domain.GoogleUser) error
	GoogleUserUpdate(ctx context.Context, key string, value any, updates map[string]any) error

	SpotifyUserInsert(ctx context.Context, spotifyUser *domain.SpotifyUser) error
	SpotifyUserUpdate(ctx context.Context, key string, value any, updates map[string]any) error

	SessionInsert(ctx context.Context, session *domain.Session) error
}

type SpotifyRepository interface {
	Me(context.Context, *oauth2.Config, *oauth2.Token) SpotifyRepository
	GetCurrentUser(context.Context) (*SpotifyPrivateUser, error)
}

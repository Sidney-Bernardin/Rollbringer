package service

import (
	"context"
	"rollbringer/pkg/domain"

	"golang.org/x/oauth2"
)

type AccountsDatabaseRepository interface {
	Close() error
	Transaction(context.Context, func(AccountsDatabaseRepository) error) error

	UserInsert(context.Context, *domain.User) error
	GoogleUserInsert(context.Context, *domain.GoogleUser) error
	SpotifyUserInsert(context.Context, *domain.SpotifyUser) error
	SessionInsert(context.Context, *domain.Session) error
}

type SpotifyRepository interface {
	Me(context.Context, *oauth2.Config, *oauth2.Token) SpotifyRepository
	GetCurrentUser(context.Context) (*SpotifyPrivateUser, error)
}

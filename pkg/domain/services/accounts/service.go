package service

import (
	"context"
	"log/slog"

	"rollbringer/pkg/domain"

	"golang.org/x/oauth2"
)

type AccountsService interface {
	domain.IService

	LoginWithGoogle(ctx context.Context, token *oauth2.Token) (*domain.User, error)
	LoginWithSpotify(ctx context.Context, oauthConfig *oauth2.Config, token *oauth2.Token) (*domain.User, error)
}

type accountsService struct {
	*domain.Service

	accountsDBRepo AccountsDatabaseRepository
	spotifyRepo    SpotifyRepository
}

func New(
	config *domain.Config,
	logger *slog.Logger,
	pubSub domain.PubSubRepository,
	accountsDBRepo AccountsDatabaseRepository,
	spotifyRepo SpotifyRepository,
) AccountsService {
	return &accountsService{
		Service: &domain.Service{
			Config: config,
			Logger: logger,
			PubSub: pubSub,
		},
		accountsDBRepo: accountsDBRepo,
		spotifyRepo:    spotifyRepo,
	}
}

func (svc *accountsService) Run(ctx context.Context) error {
	return nil
}

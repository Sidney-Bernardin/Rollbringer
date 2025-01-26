package service

import (
	"context"
	"log/slog"

	"rollbringer/pkg/domain"
)

type GamesService interface {
	domain.IService
}

type gamesService struct {
	*domain.Service

	gamesDB GamesDatabaseRepository
}

func New(config *domain.Config, logger *slog.Logger, pubSub domain.PubSubRepository, gamesDB GamesDatabaseRepository) GamesService {
	return &gamesService{
		Service: &domain.Service{
			Config: config,
			Logger: logger,
			PubSub: pubSub,
		},
		gamesDB: gamesDB,
	}
}

func (svc *gamesService) Run(ctx context.Context) error {
	return nil
}

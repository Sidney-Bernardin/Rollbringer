package users

import (
	"log/slog"

	"rollbringer/internal/config"
	"rollbringer/internal/repositories/pubsub"
	"rollbringer/internal/services"
)

type Service interface {
	services.Servicer
}

type service struct {
	*services.Service

	ps *pubsub.PubSub
}

func NewService(cfg *config.Config, logger *slog.Logger, ps *pubsub.PubSub) Service {
	return &service{
		Service: &services.Service{
			Config: cfg,
			Logger: logger,
		},
		ps: ps,
	}
}

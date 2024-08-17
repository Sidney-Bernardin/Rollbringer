package service

import (
	"log/slog"

	"rollbringer/internal/config"
	"rollbringer/internal/repositories/pubsub"
)

type UsersService interface{}

type service struct {
	cfg    *config.Config
	logger *slog.Logger

	ps *pubsub.PubSub
}

func New(cfg *config.Config, logger *slog.Logger, ps *pubsub.PubSub) UsersService {
	return &service{
		cfg:    cfg,
		logger: logger,
		ps:     ps,
	}
}

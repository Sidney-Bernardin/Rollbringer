package service

import (
	"rollbringer/pkg/repositories/database"
	"rollbringer/pkg/repositories/pubsub"

	"github.com/rs/zerolog"
)

type Service struct {
	db *database.Database
	ps *pubsub.PubSub

	logger *zerolog.Logger
}

func New(logger *zerolog.Logger, db *database.Database, ps *pubsub.PubSub) *Service {
	return &Service{
		db:     db,
		ps:     ps,
		logger: logger,
	}
}

package service

import (
	"rollbringer/pkg/repositories/database"
	"rollbringer/pkg/repositories/pubsub"
)

type Service struct {
	db *database.Database
	ps *pubsub.PubSub
}

func New(db *database.Database, ps *pubsub.PubSub) *Service {
	return &Service{
		db: db,
		ps: ps,
	}
}

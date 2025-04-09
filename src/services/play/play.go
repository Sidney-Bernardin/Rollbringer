package play

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"rollbringer/src"
	"rollbringer/src/services/play/models"

	"github.com/pkg/errors"
)

type Service interface {
	CreateRoom(context.Context, *models.Room) error
}

type service struct {
	config *src.Config

	db  Database
	bkr Broker

	canvasesUsed chan *models.EventCanvasUsed
}

func NewService(config *src.Config, db Database, bkr Broker) Service {
	return &service{
		config:       config,
		db:           db,
		bkr:          bkr,
		canvasesUsed: make(chan *models.EventCanvasUsed),
	}
}

func (svc *service) Run(ctx context.Context) error {
	return nil
}

func (svc *service) CreateRoom(ctx context.Context, room *models.Room) error {
	if len(room.Users) < 1 {
		return errors.New("room must have at least 1 user")
	}

	for i, user := range slices.Collect(maps.Values(room.Users)) {
		if len(user.Permisions) < 1 {
			return errors.New(fmt.Sprintf("room user[%d] must have at least 1 permision", i))
		}
	}

	err := svc.db.CreateRoom(ctx, room)
	return errors.Wrap(err, "database cannot create room")
}

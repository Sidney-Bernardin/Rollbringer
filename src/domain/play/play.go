package play

import (
	"context"

	"rollbringer/src"
)

const (
	ExternalErrorTypeRoomNameInvalid = "room_name_invalid"
	ExternalErrorTypeRoomNameTaken   = "room_name_taken"
)

type Service interface {
	RoomCreate(ctx context.Context, view any, cmd *ArgsRoomCreate) error
	RoomGetByID(ctx context.Context, view any, roomID string) error
}

type service struct {
	config *src.Config

	db  Database
	bkr Broker

	canvasesUsed chan *EventCanvasUsed
}

func NewService(config *src.Config, db Database, bkr Broker) Service {
	return &service{
		config:       config,
		db:           db,
		bkr:          bkr,
		canvasesUsed: make(chan *EventCanvasUsed),
	}
}

func (svc *service) Run(ctx context.Context) error {
	return nil
}

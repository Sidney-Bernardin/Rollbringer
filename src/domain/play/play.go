package play

import (
	"context"

	"rollbringer/src"
)

const (
	ExternalErrorTypeRoomNotFound    src.ExternalErrorType = "room_not_found"
	ExternalErrorTypeRoomNameInvalid src.ExternalErrorType = "room_name_invalid"
)

type Service interface {
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

package play

import (
	"context"
	"rollbringer/pkg/domain/play/commands"

	"github.com/google/uuid"
)

type (
	EventCanvasUsed struct {
		BoardID uuid.UUID
	}

	EventMovedCanvasNode struct {
		BoardID uuid.UUID
		x       int
		y       int
	}
)

type Service interface {
	RoomCreate(*commands.CreateRoom, any) error
}

type service struct {
	db  Database
	bkr Broker

	canvasesUsed chan *EventCanvasUsed
}

type Database interface {
	roomInsert(room, res any) error
	RoomGetByID(ctx context.Context, roomID commands.UUID, res any) error
}

type Broker interface {
	pubCanvaseUsed(<-chan *EventCanvasUsed)
	subCanvasesUsed(boardID commands.UUID, ch chan<- *EventCanvasUsed)

	SubMovedCanvasNodes(boardID uuid.UUID)
}

func New(db Database, bkr Broker) *service {
	return &service{
		db:           db,
		bkr:          bkr,
		canvasesUsed: make(chan *EventCanvasUsed),
	}
}

func (svc *service) Run(ctx context.Context) error {
	return nil
}

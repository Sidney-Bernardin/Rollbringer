package play

import (
	"context"

	"github.com/google/uuid"
)

type (
	Database interface {
		DatabaseCommands
		DatabaseQueries
	}

	DatabaseCommands interface{}

	DatabaseQueries interface {
		RoomGetByID(ctx context.Context, view any, roomID uuid.UUID) error
	}
)

type Broker interface {
	PubCanvaseUsed(<-chan *EventCanvasUsed)
	SubCanvasesUsed(boardID uuid.UUID, ch chan<- *EventCanvasUsed)

	SubMovedCanvasNodes(boardID uuid.UUID)
}

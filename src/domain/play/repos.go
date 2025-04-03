package play

import (
	"context"

	"rollbringer/src/domain"
)

type (
	Database interface {
		DatabaseCommands
		DatabaseQueries
	}

	DatabaseCommands interface{}

	DatabaseQueries interface {
		RoomGetByID(ctx context.Context, view any, roomID domain.UUID) error
	}
)

type Broker interface {
	PubCanvaseUsed(<-chan *EventCanvasUsed)
	SubCanvasesUsed(boardID domain.UUID, ch chan<- *EventCanvasUsed)

	SubMovedCanvasNodes(boardID domain.UUID)
}

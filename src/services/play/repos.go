package play

import (
	"context"

	"rollbringer/src"
	"rollbringer/src/services/play/models"
)

type (
	Database interface {
		DatabaseCommands
		DatabaseQueries
	}

	DatabaseCommands interface {
		CreateRoom(ctx context.Context, room *models.Room) error
	}

	DatabaseQueries interface {
		GetRoomByRoomID(ctx context.Context, roomID src.UUID) (*models.Room, error)
		GetRoomsByUserID(ctx context.Context, roomID src.UUID) ([]*models.Room, error)
	}
)

type Broker interface {
	PubCanvaseUsed(<-chan *models.EventCanvasUsed)
	SubCanvasesUsed(boardID src.UUID, ch chan<- *models.EventCanvasUsed)

	SubMovedCanvasNodes(boardID src.UUID)
}

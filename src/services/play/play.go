package play

import (
	"context"

	"rollbringer/src"
	"rollbringer/src/services"
	"rollbringer/src/services/play/models"
)

type (
	Database interface {
		CreateRoom(ctx context.Context, room *models.Room) error
		GetRoomByRoomID(ctx context.Context, roomID src.UUID) (*models.Room, error)
		GetRoomsByUserID(ctx context.Context, roomID src.UUID) ([]*models.Room, error)
		RoomExists(ctx context.Context, roomID src.UUID) (bool, error)
	}
)

type Service interface{}

type service struct {
	config *src.Config

	broker   services.Broker
	database Database
}

func NewService(config *src.Config, broker services.Broker, database Database) Service {
	return &service{config, broker, database}
}

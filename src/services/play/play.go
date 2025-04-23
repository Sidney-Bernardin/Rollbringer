package play

import (
	"context"

	"rollbringer/src"
	"rollbringer/src/services"
	"rollbringer/src/services/play/models"

	"github.com/pkg/errors"
)

const ExternalErrorTypeRoomNotFound src.ExternalErrorType = "room_not_found"

type (
	Database interface {
		BasicDatabase
		ServiceDatabase
	}

	BasicDatabase interface {
		CreateRoom(ctx context.Context, room *models.Room) error
		GetRoomsByUserID(ctx context.Context, roomID src.UUID) ([]*models.Room, error)
		RoomExists(ctx context.Context, roomID src.UUID) (bool, error)
	}

	ServiceDatabase interface {
		JoinRoom(ctx context.Context, roomUser *src.RoomUser) (*models.Room, error)
	}
)

type Service interface {
	JoinRoom(ctx context.Context, userID src.UUID, roomID src.UUID) (*models.Room, error)
}

type service struct {
	config *src.Config

	broker   services.Broker
	database Database
}

func NewService(config *src.Config, broker services.Broker, database Database) Service {
	return &service{config, broker, database}
}

func (svc *service) Chat(ctx context.Context, event *services.EventChat) error {

	// Check if the room exists.
	roomExists, err := svc.database.RoomExists(ctx, event.RoomID)
	if err != nil {
		return errors.Wrap(err, "database cannot check if room exists")
	}

	if !roomExists {
		return &src.ExternalError{
			Type: ExternalErrorTypeRoomNotFound,
			Msg:  "Cannot chat in a room that doesn't exist.",
		}
	}

	// Publish the event.
	err = svc.broker.PubChat(event)
	return errors.Wrap(err, "broker cannot publish to chat")
}

func (svc *service) JoinRoom(ctx context.Context, userID src.UUID, roomID src.UUID) (*models.Room, error) {
	room, err := svc.database.JoinRoom(ctx, &src.RoomUser{
		UserID:     userID,
		RoomID:     roomID,
		Permisions: []src.RoomUserPermision{src.RoomUserPermisionPlayer},
	})

	return room, errors.Wrap(err, "cannot join room")
}

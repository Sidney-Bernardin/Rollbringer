package play

import (
	"context"
	"rollbringer/src"
	"rollbringer/src/services"

	"github.com/pkg/errors"
)

const ExternalErrorTypeRoomNotFound src.ExternalErrorType = "room_not_found"

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

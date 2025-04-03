package play

import (
	"context"

	"github.com/pkg/errors"

	"rollbringer/src"
	"rollbringer/src/domain"
)

type RoomName string

func ParseRoomName(str string) (RoomName, error) {
	if len(str) == 0 || 30 < len(str) {
		return "", &src.ExternalError{
			Type:        ExternalErrorTypeRoomNameInvalid,
			Description: "Must be between 1 and 30 characters",
			Details:     map[string]any{"room_name": str},
		}
	}

	return RoomName(str), nil
}

/////

func (svc *service) RoomGetByID(ctx context.Context, view any, roomIDStr string) error {
	roomID, err := domain.ParseUUID(roomIDStr)
	if err != nil {
		return errors.Wrap(err, "cannot parse room-ID")
	}

	if err := svc.db.RoomGetByID(ctx, view, roomID); err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &src.ExternalError{
				Type:        ExternalErrorTypeRoomNotFound,
				Description: "Cannot find a room with the given ID",
				Details:     map[string]any{"room_id": roomID},
			}
		}

		return errors.Wrap(err, "cannot get room by ID")
	}

	return nil
}

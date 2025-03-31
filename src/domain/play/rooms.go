package play

import (
	"context"

	"github.com/google/uuid"
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
			Attrs:       map[string]any{"room_name": str},
		}
	}

	return RoomName(str), nil
}

/////

type CmdRoomCreate struct {
	OwnerID uuid.UUID
	Name    RoomName
}

type ArgsRoomCreate struct {
	OwnerID string
	Name    string
}

func (svc *service) RoomCreate(ctx context.Context, view any, args *ArgsRoomCreate) (err error) {
	var cmd CmdRoomCreate

	cmd.OwnerID, err = uuid.Parse(args.OwnerID)
	if err != nil {
		return errors.Wrap(err, "cannot parse owner-ID")
	}

	cmd.Name, err = ParseRoomName(args.OwnerID)
	if err != nil {
		return errors.Wrap(err, "cannot parse name")
	}

	if err := svc.db.RoomCreate(ctx, view, &cmd); err != nil {
		if errors.Is(err, domain.ErrEntityConflict) {
			return &src.ExternalError{
				Type:        ExternalErrorTypeRoomNameTaken,
				Description: "A room with the given name already exists.",
				Attrs:       map[string]any{"room_name": cmd.Name},
			}
		}

		return errors.Wrap(err, "cannot insert room")
	}

	return nil
}

/////

func (svc *service) RoomGetByID(ctx context.Context, view any, roomIDStr string) error {

	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		return errors.Wrap(err, "cannot parse owner-ID")
	}

	if err := svc.db.RoomGetByID(ctx, view, roomID); err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &src.ExternalError{
				Type:        src.ExternalErrorTypeEntityNotFound,
				Description: "Cannot find a room with the given ID",
				Attrs:       map[string]any{"room_id": roomID},
			}
		}

		return errors.Wrap(err, "cannot get room by ID")
	}

	return nil
}

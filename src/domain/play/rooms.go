package play

import (
	"context"

	"github.com/pkg/errors"

	"rollbringer/src/domain"
)

type RoomName string

func ParseRoomName(str string) (RoomName, error) {
	if len(str) == 0 || 30 < len(str) {
		return "", &domain.DomainError{
			Type:        DomainErrorTypeRoomNameInvalid,
			Description: "Must be between 1 and 30 characters",
			Details:     map[string]any{"room_name": str},
		}
	}

	return RoomName(str), nil
}

/////

type CmdRoomCreate struct {
	OwnerID domain.UUID
	Name    RoomName
}

type ArgsRoomCreate struct {
	OwnerID string
	Name    string
}

func (svc *service) RoomCreate(ctx context.Context, view any, args *ArgsRoomCreate) (err error) {
	var cmd CmdRoomCreate

	cmd.OwnerID, err = domain.ParseUUID(args.OwnerID)
	if err != nil {
		return errors.Wrap(err, "cannot parse owner-ID")
	}

	cmd.Name, err = ParseRoomName(args.OwnerID)
	if err != nil {
		return errors.Wrap(err, "cannot parse name")
	}

	if err := svc.db.RoomCreate(ctx, view, &cmd); err != nil {
		if errors.Is(err, domain.ErrEntityConflict) {
			return &domain.DomainError{
				Type:        DomainErrorTypeRoomNameTaken,
				Description: "A room with the given name already exists.",
				Details:     map[string]any{"room_name": cmd.Name},
			}
		}

		return errors.Wrap(err, "cannot insert room")
	}

	return nil
}

/////

func (svc *service) RoomGetByID(ctx context.Context, view any, roomIDStr string) error {
	roomID, err := domain.ParseUUID(roomIDStr)
	if err != nil {
		return errors.Wrap(err, "cannot parse owner-ID")
	}

	if err := svc.db.RoomGetByID(ctx, view, roomID); err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &domain.DomainError{
				Type:        domain.DomainErrorTypeEntityNotFound,
				Description: "Cannot find a room with the given ID",
				Details:     map[string]any{"room_id": roomID},
			}
		}

		return errors.Wrap(err, "cannot get room by ID")
	}

	return nil
}

package play

import (
	"errors"
	"rollbringer/pkg/domain/play/commands"
)

func (svc *service) RoomCreate(cmd *commands.CreateRoom, res any) error {
	if err := svc.db.RoomInsert(cmd, res); err != nil {

		if errors.Is(err, services.ErrEntityConflict) {
			return &services.DomainError{
				Type:        results.DomainErrorTypeRoomNameTaken,
				Description: "A room with the given name already exists.",
				Details:     map[string]any{"room_name": cmd.Name},
			}
		}

		return errors.Join(err, errors.New("cannot insert room"))
	}

	return nil
}

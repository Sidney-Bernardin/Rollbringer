package play

import (
	"context"
	"errors"

	"rollbringer/server/domain"
	"rollbringer/server/domain/play/commands"
	"rollbringer/server/domain/play/results"
)

func (svc *service) RoomCreate(ctx context.Context, cmd *commands.RoomCreate, res any) error {
	if err := svc.db.RoomInsert(ctx, cmd, res); err != nil {
		if errors.Is(err, domain.ErrEntityConflict) {
			return &domain.DomainError{
				Type:        results.DomainErrorTypeRoomNameTaken,
				Description: "A room with the given name already exists.",
				Details:     map[string]any{"room_name": cmd.Name},
			}
		}

		return errors.Join(err, errors.New("cannot insert room"))
	}

	return nil
}

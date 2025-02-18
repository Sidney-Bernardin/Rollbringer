package service

import (
	"context"
	"log/slog"
	"rollbringer/pkg/domain"
)

func (svc *playService) subPlay(ctx context.Context, e *domain.Event) *domain.Event {
	switch p := e.Payload.(type) {

	case *domain.GetRoomRequest:
		room, err := svc.GetRoom(ctx, p.RoomID)
		if err != nil {
			return &domain.Event{
				Operation: domain.OperationError,
				Payload:   domain.HandleError(ctx, svc.Logger, slog.LevelError, domain.Wrap(err, "cannot get room", nil)),
			}
		} else {
			return &domain.Event{
				Operation: domain.OperationRoom,
				Payload:   room,
			}
		}

	case *domain.GetRoomsRequest:
		rooms, err := svc.GetRooms(ctx, p.OwnerID)
		if err != nil {
			return &domain.Event{
				Operation: domain.OperationError,
				Payload:   domain.HandleError(ctx, svc.Logger, slog.LevelError, domain.Wrap(err, "cannot get rooms", nil)),
			}
		} else {
			return &domain.Event{
				Operation: domain.OperationRooms,
				Payload:   rooms,
			}
		}

	default:
		return nil
	}
}

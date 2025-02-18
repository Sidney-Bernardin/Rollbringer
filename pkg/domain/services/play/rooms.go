package service

import (
	"context"

	"github.com/google/uuid"

	"rollbringer/pkg/domain"
)

func (svc *playService) CreateRoom(ctx context.Context, session *domain.Session, room *domain.Room) error {
	err := svc.playDBRepo.Transaction(ctx, func(tx PlayDatabaseRepository) error {
		room.OwnerID = session.UserID

		if err := tx.RoomInsert(ctx, room); err != nil {
			return domain.Wrap(err, "cannot insert room", nil)
		}

		board := &domain.Board{
			RoomID: room.ID,
			Name:   "Test Board",
		}

		if err := tx.BoardInsert(ctx, board); err != nil {
			return domain.Wrap(err, "cannot insert board", nil)
		}

		return nil
	})

	return domain.Wrap(err, "cannot do transaction", nil)
}

func (svc *playService) GetRoom(ctx context.Context, roomID uuid.UUID) (*domain.Room, error) {
	room, err := svc.playDBRepo.RoomGet(ctx, "id", roomID)
	return room, domain.Wrap(err, "cannot get room", nil)
}

func (svc *playService) GetRooms(ctx context.Context, ownerID uuid.UUID) ([]*domain.Room, error) {
	rooms, err := svc.playDBRepo.RoomsGet(ctx, "owner_id", ownerID)
	return rooms, domain.Wrap(err, "cannot get rooms", nil)
}

func (svc *playService) DeleteRoom(ctx context.Context, session *domain.Session, roomID uuid.UUID) error {
	err := svc.playDBRepo.RoomDelete(ctx, roomID, session.UserID)
	return domain.Wrap(err, "cannot delete room", nil)
}

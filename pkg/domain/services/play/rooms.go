package service

import (
	"context"

	"github.com/google/uuid"

	"rollbringer/pkg/domain"
)

func (svc *playService) CreateRoom(ctx context.Context, session *domain.Session, room *domain.Room) error {
	room.OwnerID = session.UserID

	err := svc.playDBRepo.RoomInsert(ctx, room)
	return domain.Wrap(err, "cannot insert room", nil)
}

func (svc *playService) GetRoom(ctx context.Context, roomID uuid.UUID) (*domain.Room, error) {

	room, err := svc.playDBRepo.RoomGet(ctx, "id", roomID)
	if err != nil {
		return nil, domain.Wrap(err, "cannot get room", nil)
	}

	room.Boards, err = svc.playDBRepo.BoardsGet(ctx, domain.BoardViewListItem, "room_id", room.ID)
	if err != nil {
		return nil, domain.Wrap(err, "cannot get boards", nil)
	}

	return room, nil
}

func (svc *playService) GetRooms(ctx context.Context, ownerID uuid.UUID) ([]*domain.Room, error) {
	rooms, err := svc.playDBRepo.RoomsGet(ctx, "owner_id", ownerID)
	return rooms, domain.Wrap(err, "cannot get rooms", nil)
}

func (svc *playService) DeleteRoom(ctx context.Context, session *domain.Session, roomID uuid.UUID) error {
	err := svc.playDBRepo.RoomDelete(ctx, roomID, session.UserID)
	return domain.Wrap(err, "cannot delete room", nil)
}

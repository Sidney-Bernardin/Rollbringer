package database

import (
	"context"
	"fmt"

	"rollbringer/pkg/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const qRoomInsert = ` 
WITH inserted_room AS (
	INSERT INTO play.rooms (name, owner_id)
	VALUES ($1, $2)
	RETURNING *
)
SELECT * FROM inserted_room`

func (repo *playDatabaseRepository) RoomInsert(ctx context.Context, room *domain.Room) error {
	err := repo.Insert(ctx, room, qRoomInsert,
		room.Name, room.OwnerID)
	return domain.Wrap(err, "cannot insert room", nil)
}

/////

const qRoomsGet = ` 
SELECT * FROM play.rooms WHERE rooms.%s = $1`

func (repo *playDatabaseRepository) RoomGet(ctx context.Context, key string, value any) (*domain.Room, error) {
	room := &domain.Room{}
	if err := repo.GetOne(ctx, room, fmt.Sprintf(qRoomsGet, key), value); err != nil {
		return nil, domain.Wrap(err, "cannot select room", nil)
	}
	return room, nil
}

func (repo *playDatabaseRepository) RoomsGet(ctx context.Context, key string, value any) ([]*domain.Room, error) {
	rooms := []*domain.Room{}
	if err := sqlx.SelectContext(ctx, repo.TX, &rooms, fmt.Sprintf(qRoomsGet, key), value); err != nil {
		return nil, domain.Wrap(err, "cannot select rooms", nil)
	}
	return rooms, nil
}

/////

const qRoomsDelete = ` 
DELETE FROM play.rooms WHERE rooms.id = $1 AND rooms.owner_id = $2`

func (repo *playDatabaseRepository) RoomDelete(ctx context.Context, roomID uuid.UUID, ownerID uuid.UUID) error {
	_, err := repo.TX.ExecContext(ctx, qRoomsDelete, roomID, ownerID)
	return domain.Wrap(err, "cannot delete room", nil)
}

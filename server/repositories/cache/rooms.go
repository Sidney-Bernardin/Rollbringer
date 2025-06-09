package cache

import (
	"context"
	"fmt"

	"github.com/Sidney-Bernardin/Rollbringer/server"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/sql/queries"

	"github.com/pkg/errors"
)

func (cache *Cache) SetRoom(ctx context.Context, room *queries.SelectRoomRow) error {
	key := fmt.Sprintf("room:%s", room.ID)
	err := cache.client.JSONSet(ctx, key, "$", room).Err()
	return errors.WithStack(err)
}

func (cache *Cache) GetRoom(ctx context.Context, roomID server.UUID) (*queries.SelectRoomRow, error) {
	key := fmt.Sprintf("room:%s", roomID)
	room, err := jsonGet[*queries.SelectRoomRow](ctx, cache.client, key, ".")
	return room, errors.WithStack(err)
}

func (cache *Cache) DeleteRoom(ctx context.Context, roomID server.UUID) error {
	query := fmt.Sprintf(`@room_id:{%s} | @user_rooms_ids:{%s}`, roomID.Escape(), roomID.Escape())
	err := ftDelete(ctx, cache.client, "room_index", query)
	return errors.WithStack(err)
}

/////

func (cache *Cache) SetUserRooms(ctx context.Context, userID server.UUID, rooms []*queries.SelectUserRoomsRow) error {
	key := fmt.Sprintf("user_rooms:%s", userID)
	err := cache.client.JSONSet(ctx, key, "$", rooms).Err()
	return errors.WithStack(err)
}

func (cache *Cache) GetUserRooms(ctx context.Context, userID server.UUID) ([]*queries.SelectUserRoomsRow, error) {
	key := fmt.Sprintf("user_rooms:%s", userID)
	rooms, err := jsonGet[[]*queries.SelectUserRoomsRow](ctx, cache.client, key, ".")
	return rooms, errors.WithStack(err)
}

func (cache *Cache) DeleteUserRooms(ctx context.Context, userID server.UUID) error {
	key := fmt.Sprintf("user_rooms:%s", userID)
	err := cache.client.Del(ctx, key).Err()
	return errors.WithStack(err)
}

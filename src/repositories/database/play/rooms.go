package play

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"rollbringer/src"
	"rollbringer/src/services/play/models"
)

type sqlRoom struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

func (r *sqlRoom) domain() *models.Room {
	return &models.Room{
		ID:    src.UUID(r.ID),
		Name:  models.RoomName(r.Name),
		Users: map[src.UUID]*models.RoomUser{},
	}
}

type sqlRoomUser struct {
	RoomID     uuid.UUID `db:"room_id"`
	UserID     uuid.UUID `db:"user_id"`
	Permisions []string  `db:"permisions"`
}

func (ru *sqlRoomUser) domain() *models.RoomUser {
	roomUser := &models.RoomUser{
		UserID:     src.UUID(ru.UserID),
		Permisions: make([]models.RoomUserPermision, len(ru.Permisions)),
	}

	for i, permision := range ru.Permisions {
		roomUser.Permisions[i] = models.RoomUserPermision(permision)
	}

	return roomUser
}

/////

const qInsertRoomAndRoomUser = `
	WITH 
		inserted_room AS (
			INSERT INTO play.rooms (id, name)
			VALUES ($2, $3)
			RETURNING *
		),
		inserted_room_user_permisions AS (
			INSERT INTO play.room_user_permissions (room_id, user_id)
			VALUES ($2, $1)
		)`

func (db *playDatabase) CreateRoom(ctx context.Context, room *models.Room) error {
	creator := slices.Collect(maps.Values(room.Users))[0]

	var r *sqlRoom
	err := db.Insert(ctx, r, qInsertRoomAndRoomUser,
		creator.UserID, room.ID, room.Name)
	if err != nil {
		return errors.Wrap(err, "cannot create room")
	}

	*room = *r.domain()
	room.Users[creator.UserID] = creator
	return nil
}

/////

const (
	qSelectRoomByRoomID = `
		SELECT rooms.id, rooms.name, room_users.permisions
		FROM play.rooms
		WHERE rooms.id = $1`

	qSelectRoomUsersByRoomID = `
		SELECT room_users.user_id, room_users.permisions
		FROM play.room_users
		WHERE room_users.room_id = $1`
)

func (db *playDatabase) GetRoomByRoomID(ctx context.Context, roomID src.UUID) (*models.Room, error) {

	var r *sqlRoom
	if err := db.SelectOne(ctx, r, qSelectRoomByRoomID, roomID); err != nil {
		return nil, errors.Wrap(err, "cannot select room")
	}

	var uu []sqlRoomUser
	if err := db.SelectMany(ctx, uu, qSelectRoomUsersByRoomID, roomID); err != nil {
		return nil, errors.Wrap(err, "cannot select room-users")
	}

	room := r.domain()
	for _, u := range uu {
		room.Users[src.UUID(u.UserID)] = u.domain()
	}

	return room, nil
}

/////

const qSelectRoomsByUserID = `
	SELECT rooms.id, rooms.name, room_users.permisions
	FROM play.rooms
	WHERE EXISTS (
		SELECT 1 FROM play.room_users WHERE room_users.room_id = rooms.id AND room_users.user_id = $1
	)`

const qSelectRoomUsersByRoomIDOrUserID = `
	SELECT room_users.room_id, room_users.user_id, room_users.permisions
	FROM play.room_users
	WHERE room_users.room_id IN (%s)`

func (db *playDatabase) GetRoomsByUserID(ctx context.Context, userID src.UUID) ([]*models.Room, error) {

	var rr []*sqlRoom
	if err := db.SelectMany(ctx, rr, qSelectRoomsByUserID, userID); err != nil {
		return nil, errors.Wrap(err, "cannot create room")
	}

	roomIDs := make([]string, len(rr))
	for i, r := range rr {
		roomIDs[i] = r.ID.String()
	}

	var uu []sqlRoomUser
	if err := db.SelectMany(ctx, uu, fmt.Sprintf(qSelectRoomUsersByRoomIDOrUserID, strings.Join(roomIDs, ","))); err != nil {
		return nil, errors.Wrap(err, "cannot select room-users")
	}

	rooms := make([]*models.Room, len(rr))
	for i, r := range rr {
		rooms[i] = r.domain()
		for _, u := range uu {
			if u.RoomID == r.ID {
				rooms[i].Users[src.UUID(u.UserID)] = u.domain()
			}
		}
	}

	return rooms, nil
}

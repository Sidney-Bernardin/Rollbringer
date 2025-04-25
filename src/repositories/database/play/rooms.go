package play

import (
	"context"
	"maps"
	"slices"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"rollbringer/src"
	"rollbringer/src/repositories/database"
	"rollbringer/src/services/play/models"
)

type roomRow struct {
	ID   pgtype.UUID `db:"rooms.id"`
	Name string      `db:"rooms.name"`

	UserIDs        []pgtype.UUID             `db:"room_user_permisions.user_ids"`
	UserPermisions [][]src.RoomUserPermision `db:"room_user_permisions.permisions"`
}

func (r *roomRow) Domain() *models.Room {
	if r == nil {
		return nil
	}

	users := make(map[src.UUID][]src.RoomUserPermision, len(r.UserIDs))
	for i, userID := range r.UserIDs {
		users[src.UUID(userID.Bytes)] = r.UserPermisions[i]
	}

	return &models.Room{
		ID:             src.UUID(r.ID.Bytes),
		Name:           models.RoomName(r.Name),
		UserPermisions: users,
	}
}

func (db *playDatabase) CreateRoom(ctx context.Context, room *models.Room) error {
	creatorID := slices.Collect(maps.Keys(room.UserPermisions))[0]

	err := database.Insert(ctx, db.Tx, `
		WITH inserted_room AS (
			INSERT INTO play.rooms (id, name)
			VALUES ($1, $2)
		)
		INSERT INTO room_user_permisions (room_id, user_id, permisions)
		VALUES ($1, $3, $4)
	`, room.ID, room.Name, creatorID, pq.Array(room.UserPermisions[creatorID]))

	return errors.Wrap(err, "cannot create room")
}

func (db *playDatabase) JoinRoom(ctx context.Context, userID, roomID src.UUID, permisions ...src.RoomUserPermision) (*models.Room, error) {
	row, err := database.Get[roomRow](ctx, db.Tx, `
		WITH inserted_room_user_permision AS (
			INSERT INTO room_user_permisions (room_id, user_id, permisions)
			VALUES ($1, $2, $3)
			ON CONFLICT (room_id, user_id) DO NOTHING
		)
		SELECT
			rooms.id AS "rooms.id",
			rooms.name AS "rooms.name",
			json_agg(room_user_permisions.user_id) AS "room_user_permisions.user_ids",
			json_agg(room_user_permisions.permisions) AS "room_user_permisions.permisions"
		FROM play.rooms
		LEFT JOIN room_user_permisions ON rooms.id = room_user_permisions.room_id
		WHERE rooms.id = $1
		GROUP BY rooms.id
	`, roomID, userID, pq.Array(permisions))

	return row.Domain(), errors.Wrap(err, "cannot run query")
}

func (db *playDatabase) GetRoomsByUserID(ctx context.Context, userID src.UUID) ([]*models.Room, error) {
	rooms, err := database.Gets[roomRow](ctx, db.Tx, `
		SELECT
			rooms.id AS "rooms.id",
			rooms.name AS "rooms.name",
			json_agg(room_user_permisions.user_id) AS "room_user_permisions.user_ids",
			json_agg(room_user_permisions.permisions) AS "room_user_permisions.permisions"
		FROM play.rooms
		LEFT JOIN room_user_permisions ON rooms.id = room_user_permisions.room_id
		WHERE room_user_permisions.user_id = $1
		GROUP BY rooms.id
	`, userID)

	return database.Domains(rooms), errors.Wrap(err, "cannot select rooms by user-ID")
}

func (db *playDatabase) RoomExists(ctx context.Context, roomID src.UUID) (exists bool, err error) {
	err = db.Tx.QueryRow(ctx, `
		SELECT EXISTS (SELECT * FROM play.rooms WHERE id = $1)
	`, roomID).Scan(&exists)

	return exists, errors.Wrap(err, "cannot select room by room-ID")
}

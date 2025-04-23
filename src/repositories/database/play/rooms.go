package play

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"rollbringer/src"
	"rollbringer/src/repositories/database"
	"rollbringer/src/services/play/models"
)

type roomRow struct {
	ID             pgtype.UUID               `db:"rooms.id"`
	Name           string                    `db:"rooms.name"`
	UserIDs        []pgtype.UUID             `db:"room_users.user_ids"`
	UserPermisions [][]src.RoomUserPermision `db:"room_users.permisions"`
}

func (r *roomRow) Domain() *models.Room {
	if r == nil {
		return nil
	}

	users := make([]*src.RoomUser, len(r.UserIDs))
	for i, userID := range r.UserIDs {
		users[i] = &src.RoomUser{
			UserID:     src.UUID(userID.Bytes),
			Permisions: r.UserPermisions[i],
		}
	}

	return &models.Room{
		ID:    src.UUID(r.ID.Bytes),
		Name:  models.RoomName(r.Name),
		Users: users,
	}
}

func (db *playDatabase) CreateRoom(ctx context.Context, room *models.Room) error {
	err := database.Insert(ctx, db.Tx, `
		WITH inserted_room AS (
			INSERT INTO play.rooms (id, name)
			VALUES ($1, $2)
		)
		INSERT INTO room_users (room_id, user_id, permisions)
		VALUES ($1, $3, $4)
	`,
		room.ID, room.Name, room.Users[0].UserID, pq.Array(room.Users[0].Permisions))

	return errors.Wrap(err, "cannot create room")
}

func (db *playDatabase) JoinRoom(ctx context.Context, roomUser *src.RoomUser) (room *models.Room, err error) {
	row, err := database.Get[roomRow](ctx, db.Tx, `
		WITH inserted_room_user AS (
			INSERT INTO room_users (room_id, user_id, permisions)
			VALUES ($1, $2, $3)
			ON CONFLICT (room_id, user_id) DO NOTHING
		)
		SELECT
			rooms.id AS "rooms.id",
			rooms.name AS "rooms.name",
			json_agg(room_users.user_id) AS "room_users.user_ids",
			json_agg(room_users.permisions) AS "room_users.permisions"
		FROM play.rooms
		LEFT JOIN room_users ON rooms.id = room_users.room_id
		WHERE rooms.id = $1
		GROUP BY rooms.id
	`, roomUser.RoomID, roomUser.UserID, pq.Array(roomUser.Permisions))

	return row.Domain(), errors.Wrap(err, "cannot run query")
}

func (db *playDatabase) GetRoomsByUserID(ctx context.Context, userID src.UUID) ([]*models.Room, error) {
	rooms, err := database.Gets[roomRow](ctx, db.Tx, `
		SELECT
			rooms.id AS "rooms.id",
			rooms.name AS "rooms.name",
			json_agg(room_users.user_id) AS "room_users.user_ids",
			json_agg(room_users.permisions) AS "room_users.permisions"
		FROM play.rooms
		LEFT JOIN room_users ON rooms.id = room_users.room_id
		WHERE room_users.user_id = $1
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

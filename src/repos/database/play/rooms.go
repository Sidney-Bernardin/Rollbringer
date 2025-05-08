package play

import (
	"context"
	"maps"
	"slices"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"rollbringer/src/domain/services/play"
	"rollbringer/src/repos/database"
)

type room struct {
	ID   pgtype.UUID `db:"rooms.id"`
	Name string      `db:"rooms.name"`

	UserIDs    []uuid.UUID                `db:"room_users.user_ids"`
	Permisions [][]play.RoomUserPermision `db:"room_users.permisions"`
}

func (r room) Model() *play.Room {
	permisions := make(map[uuid.UUID][]play.RoomUserPermision, len(r.UserIDs))
	for i, userID := range r.UserIDs {
		permisions[userID] = r.Permisions[i]
	}

	return &play.Room{
		ID:             uuid.UUID(r.ID.Bytes),
		Name:           play.RoomName(r.Name),
		UserPermisions: permisions,
	}
}

func (db *playDatabase) CreateRoom(ctx context.Context, room *play.Room) error {
	err := db.Transaction(ctx, func(tx pgx.Tx) error {

		// Insert the room.
		err := database.Insert(ctx, db.Pool,
			`
				INSERT INTO play.rooms (id, name)
				VALUES ($1, $2)
			`,
			room.ID, room.Name)

		if err != nil {
			return errors.Wrap(err, "cannot insert room")
		}

		// Copy the user-permisions into the play.room_users table.
		userIDs := slices.Collect(maps.Keys(room.UserPermisions))
		_, err = tx.CopyFrom(ctx,
			pgx.Identifier{"play", "room_users"},
			[]string{"room_id", "user_id", "permisions"},
			pgx.CopyFromSlice(
				len(userIDs),
				func(i int) ([]any, error) {
					uID := userIDs[i]
					userPermisions := pgtype.FlatArray[string]([]string{})
					for _, p := range room.UserPermisions[uID] {
						userPermisions = append(userPermisions, string(p))
					}

					return []any{room.ID, uID, userPermisions}, nil
				}))

		return errors.Wrap(err, "cannot insert room_user")
	})

	return errors.Wrap(err, "transaction failed")
}

func (db *playDatabase) JoinRoom(ctx context.Context, userID, roomID uuid.UUID, permisions ...play.RoomUserPermision) (*play.Room, bool, error) {
	_, model, err := database.Get[room](ctx, db.Pool,
		`
			WITH selection AS (
				SELECT
					rooms.id AS "rooms.id",
					rooms.name AS "rooms.name",
					json_agg(room_users.user_id) AS "room_users.user_ids",
					json_agg(room_users.permisions) AS "room_users.permisions"
				FROM play.rooms
				LEFT JOIN play.room_users ON rooms.id = room_users.room_id
				WHERE rooms.id = $1
				GROUP BY rooms.id
			), inserted AS (
				INSERT INTO play.room_users (room_id, user_id, permisions)
				VALUES ($1, $2, $3)
				ON CONFLICT (room_id, user_id) DO NOTHING
			) SELECT * FROM selection
		`,
		roomID, userID, pq.Array(permisions))

	if err != nil {
		return nil, false, errors.Wrap(err, "cannot run query")
	}

	if _, ok := model.UserPermisions[userID]; ok {
		return model, false, nil
	}

	model.UserPermisions[userID] = permisions
	return model, true, nil
}

func (db *playDatabase) GetRoomsByUserID(ctx context.Context, userID uuid.UUID) ([]*play.Room, error) {
	_, models, err := database.Gets[room](ctx, db.Pool,
		`
			SELECT
				rooms.id AS "rooms.id",
				rooms.name AS "rooms.name",
				json_agg(room_users.user_id) AS "room_users.user_ids",
				json_agg(room_users.permisions) AS "room_users.permisions"
			FROM play.rooms
			LEFT JOIN play.room_users ON rooms.id = room_users.room_id
			WHERE room_users.user_id = $1
			GROUP BY rooms.id
		`,
		userID)

	return models, errors.Wrap(err, "cannot select rooms by user-ID")
}

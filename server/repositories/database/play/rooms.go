package play

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"rollbringer/server/domain/play/commands"
	"rollbringer/server/domain/play/results"
	"rollbringer/server/repositories/database"
)

type room struct {
	ID      uuid.UUID `db:"id"`
	OwnerID uuid.UUID `db:"owner_id"`
	Name    string    `db:"name"`
}

func (db *playDatabase) roomQuery(ctx context.Context, rowFunc database.RowFunc, view any, query string, args ...any) error {

	var columns, joins string
	switch view.(type) {
	default:
		columns = `rooms.id, rooms.name`
	case results.RoomListItem:
		columns = `rooms.id, rooms.owner_id, rooms.name`
	}

	var r *room
	if err := rowFunc(ctx, r, fmt.Sprintf(query, columns, joins), args...); err != nil {
		return err
	}

	switch v := view.(type) {
	case *results.RoomInfo:
		v.RoomID = r.ID
		v.OwnerID = r.OwnerID
		v.Name = r.Name
	case *results.RoomListItem:
		v.RoomID = r.ID
		v.Name = r.Name
	}

	return nil
}

/////

const qRoomInsert = `
	WITH inserted_room AS (
		INSERT INTO play.rooms (host_id, name)
		VALUES ($1, $2)
		RETURNING *
	)
	SELECT %s FROM inserted_room %s`

func (db *playDatabase) RoomCreate(ctx context.Context, cmd *commands.RoomCreate, view any) error {
	return db.roomQuery(ctx, db.RowInsert, view, qRoomInsert, cmd.HostID, cmd.Name)
}

/////

const qRoomGetByID = `SELECT %s FROM rooms %s`

func (db *playDatabase) RoomGetByID(ctx context.Context, roomID commands.UUID, view any) error {
	return db.roomQuery(ctx, db.RowGet, view, qRoomGetByID, roomID)
}

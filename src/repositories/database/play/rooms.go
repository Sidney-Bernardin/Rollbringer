package play

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"rollbringer/src/domain"
	"rollbringer/src/domain/play"
	"rollbringer/src/repositories/database"
)

type room struct {
	ID      uuid.UUID `db:"id"`
	OwnerID uuid.UUID `db:"owner_id"`
	Name    string    `db:"name"`
}

func (db *playDatabase) roomQuery(ctx context.Context, crudFunc database.CRUDFunc, view any, query string, args ...any) error {

	var columns, joins string
	switch view.(type) {
	default:
		columns = `rooms.id, rooms.name`
	case play.RoomListItem:
		columns = `rooms.id, rooms.owner_id, rooms.name`
	}

	var r room
	if err := crudFunc(ctx, &r, fmt.Sprintf(query, columns, joins), args...); err != nil {
		return err
	}

	switch v := view.(type) {
	case *play.RoomInfo:
		v.RoomID = r.ID.String()
		v.OwnerID = r.OwnerID.String()
		v.Name = r.Name
	case *play.RoomListItem:
		v.RoomID = r.ID.String()
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

func (db *playDatabase) RoomCreate(ctx context.Context, view any, cmd *play.CmdRoomCreate) error {
	return db.roomQuery(ctx, db.CRUDInsert, view, qRoomInsert, cmd.OwnerID, cmd.Name)
}

/////

const qRoomGetByID = `SELECT %s FROM play.rooms %s WHERE id = $1`

func (db *playDatabase) RoomGetByID(ctx context.Context, view any, roomID domain.UUID) error {
	return db.roomQuery(ctx, db.CRUDGet, view, qRoomGetByID, roomID)
}

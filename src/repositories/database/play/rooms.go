package play

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"rollbringer/src/domain"
	"rollbringer/src/domain/play"
	"rollbringer/src/repositories/database"
)

type room struct {
	ID      domain.UUID `db:"id"`
	OwnerID domain.UUID `db:"owner_id"`
	Name    string      `db:"name"`
}

const (
	qRoomSelectByID = `SELECT %s FROM play.rooms %s WHERE id = $1`
)

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

func (db *playDatabase) RoomGetByID(ctx context.Context, view any, roomID domain.UUID) error {
	err := db.roomQuery(ctx, db.CRUDGet, view, qRoomSelectByID, roomID)
	return errors.Wrap(err, "cannot get room by ID")
}

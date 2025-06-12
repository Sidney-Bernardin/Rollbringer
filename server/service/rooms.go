package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/Sidney-Bernardin/Rollbringer/server"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/cache"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/sql"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/sql/queries"

	"github.com/pkg/errors"
)

func (svc *Service) CreateRoom(ctx context.Context, creatorID server.UUID, name string) (*queries.InsertRoomRow, error) {

	if len(name) < 3 || 32 < len(name) {
		return nil, &server.UserError{
			Type:    server.UserErrorTypeRoomNameInvalid,
			Message: "Username must be between 3 and 32 characters long.",
		}
	}

	room, err := svc.SQL.InsertRoom(ctx, &queries.InsertRoomParams{
		UserID: creatorID,
		RoomID: server.NewRandomUUID(),
		Name:   name,
		Permisions: []queries.UserRoomPermision{
			queries.UserRoomPermisionOWNER,
			queries.UserRoomPermisionGAMEMASTER,
		}})
	if err != nil {
		return nil, errors.Wrap(err, "cannot insert room")
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := svc.Cache.DeleteUserRooms(ctx, creatorID); err != nil {
			svc.Log.Log(ctx, slog.LevelWarn, "Cache cannot delete user's rooms",
				"err", err.Error(),
				"deleter_id", creatorID)
		}
	}()

	return room, nil
}

func (svc *Service) GetRoom(ctx context.Context, roomID server.UUID) (*queries.SelectRoomRow, error) {

	room, err := svc.Cache.GetRoom(ctx, roomID)
	if !errors.Is(err, cache.ErrNotFound) {
		return room, errors.Wrap(err, "cannot get room from Cache")
	}

	room, err = svc.SQL.SelectRoom(ctx, roomID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &server.UserError{
				Type:    server.UserErrorTypeRoomNotFound,
				Message: "Cannot find a room with that ID",
			}
		}

		return nil, errors.Wrap(err, "cannot get room from SQL")
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := svc.Cache.SetRoom(ctx, room); err != nil {
			svc.Log.Log(ctx, slog.LevelWarn, "Cache cannot set room",
				"err", err.Error(),
				"room_id", roomID)
		}
	}()

	return room, nil
}

func (svc *Service) GetUserRooms(ctx context.Context, userID server.UUID) ([]*queries.SelectUserRoomsRow, error) {

	rooms, err := svc.Cache.GetUserRooms(ctx, userID)
	if !errors.Is(err, cache.ErrNotFound) {
		return rooms, errors.Wrap(err, "cannot get rooms from Cache")
	}

	rooms, err = svc.SQL.SelectUserRooms(ctx, userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "cannot get rooms from SQL")
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := svc.Cache.SetUserRooms(ctx, userID, rooms); err != nil {
			svc.Log.Log(ctx, slog.LevelWarn, "Cache cannot set user's rooms",
				"err", err.Error(),
				"setter_id", userID)
		}
	}()

	return rooms, nil
}

func (svc *Service) DeleteRoom(ctx context.Context, userID, roomID server.UUID) error {
	err := svc.SQL.DeleteRoom(ctx, &queries.DeleteRoomParams{
		RoomID: roomID,
		UserID: userID})
	if err != nil {
		return errors.WithStack(err)
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := svc.Cache.DeleteRoom(ctx, roomID); err != nil {
			svc.Log.Log(ctx, slog.LevelWarn, "Cache cannot delete room",
				"err", err.Error(),
				"room_id", roomID,
				"deleter_id", userID)
		}
	}()

	return nil
}

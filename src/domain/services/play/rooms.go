package play

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"rollbringer/src/domain"
)

const (
	ExternalErrorTypeInvalidRoomName domain.ExternalErrorType = "invalids-rooms-name"
)

type Room struct {
	ID             uuid.UUID
	Name           RoomName
	UserPermisions map[uuid.UUID][]RoomUserPermision
}

type RoomName string

func ParseRoomName(str string) (RoomName, error) {
	if len(str) == 0 || 30 < len(str) {
		return "", &domain.ExternalError{
			Type: ExternalErrorTypeInvalidRoomName,
			Msg:  "Must be between 1 and 30 characters",
			Details: map[string]any{
				"room_name": str,
			},
		}
	}

	return RoomName(str), nil
}

type RoomUserPermision string

const (
	RoomUserPermisionOwner      RoomUserPermision = "OWNER"
	RoomUserPermisionGameMaster RoomUserPermision = "GAME_MASTER"
	RoomUserPermisionPlayer     RoomUserPermision = "PLAYER"
)

type CreateRoomOpts struct {
	Name string `json:"name"`
}

func (svc *service) CreateRoom(ctx context.Context, creatorID uuid.UUID, opts *CreateRoomOpts) (*Room, error) {

	name, err := ParseRoomName(opts.Name)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse name")
	}

	room := &Room{
		ID:   uuid.New(),
		Name: name,
		UserPermisions: map[uuid.UUID][]RoomUserPermision{
			creatorID: {RoomUserPermisionOwner, RoomUserPermisionGameMaster},
		},
	}

	if err := svc.db.CreateRoom(ctx, room); err != nil {
		return nil, errors.Wrap(err, "database cannot create room")
	}

	return room, nil
}

func (svc *service) JoinRoom(ctx context.Context, roomID uuid.UUID, newcomer *domain.PublicUser) (*Room, error) {

	room, newlyJoined, err := svc.db.JoinRoom(ctx, newcomer.UserID, roomID, RoomUserPermisionPlayer)
	if err != nil {
		return nil, errors.Wrap(err, "cannot join room")
	}

	if newlyJoined {
		svc.broker.Pub(ctx, &domain.EventRoomJoined{
			RoomID:   roomID,
			Newcomer: *newcomer,
		})
	}

	return room, nil
}

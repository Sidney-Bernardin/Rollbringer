package models

import (
	"github.com/pkg/errors"

	"rollbringer/src"
)

const (
	ExternalErrorTypeInvalidRoomName src.ExternalErrorType = "invalid_room_name"
)

type Room struct {
	ID             src.UUID                             `json:"id"`
	Name           RoomName                             `json:"name"`
	UserPermisions map[src.UUID][]src.RoomUserPermision `json:"user_permisions"`
}

func NewRoom(creatorID src.UUID, name string) (*Room, error) {
	roomName, err := ParseRoomName(name)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse room-name")
	}

	return &Room{
		ID:   src.NewUUID(),
		Name: roomName,
		UserPermisions: map[src.UUID][]src.RoomUserPermision{
			creatorID: {src.RoomUserPermisionOwner, src.RoomUserPermisionGameMaster},
		},
	}, nil
}

type RoomName string

func ParseRoomName(str string) (RoomName, error) {
	if len(str) == 0 || 30 < len(str) {
		return "", &src.ExternalError{
			Type: ExternalErrorTypeInvalidRoomName,
			Msg:  "Must be between 1 and 30 characters",
			Details: map[string]any{
				"room_name": str,
			},
		}
	}

	return RoomName(str), nil
}

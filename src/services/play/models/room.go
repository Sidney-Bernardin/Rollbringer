package models

import (
	"rollbringer/src"
)

type Room struct {
	ID    src.UUID               `json:"id"`
	Name  RoomName               `json:"name"`
	Users map[src.UUID]*RoomUser `json:"users"`
}

type RoomUser struct {
	UserID     src.UUID            `json:"user_id"`
	Permisions []RoomUserPermision `json:"permisions"`
}

type RoomUserPermision string

const (
	RoomUserPermisionOwner      RoomUserPermision = "OWNER"
	RoomUserPermisionGameMaster RoomUserPermision = "GAME_MASTER"
	RoomUserPermisionPlayer     RoomUserPermision = "PLAYER"
)

type RoomName string

const ExternalErrorTypeInvalidRoomName src.ExternalErrorType = "INVALID_ROOM_NAME"

func ParseRoomName(str string) (RoomName, error) {
	if len(str) == 0 || 30 < len(str) {
		return "", &src.ExternalError{
			Type:        ExternalErrorTypeInvalidRoomName,
			Description: "Must be between 1 and 30 characters",
			Details: map[string]any{
				"room_name": str,
			},
		}
	}

	return RoomName(str), nil
}

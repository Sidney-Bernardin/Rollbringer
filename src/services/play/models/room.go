package models

import (
	"rollbringer/src"
)

type Room struct {
	ID    src.UUID        `json:"id"`
	Name  RoomName        `json:"name"`
	Users []*src.RoomUser `json:"users"`
}

type RoomName string

const ExternalErrorTypeInvalidRoomName src.ExternalErrorType = "INVALID_ROOM_NAME"

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

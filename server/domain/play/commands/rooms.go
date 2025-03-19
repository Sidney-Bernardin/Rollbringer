package commands

import (
	"rollbringer/server/domain"
	"rollbringer/server/domain/play/results"
)

type RoomName string

func NewRoomName(name string) (RoomName, error) {
	if len(name) == 0 || 30 < len(name) {
		return "", &domain.DomainError{
			Type:        results.DomainErrorTypeRoomNameInvalid,
			Description: "Must be between 1 and 30 characters",
			Details:     map[string]any{"room_name": name},
		}
	}

	return RoomName(name), nil
}

type RoomCreate struct {
	HostID UUID
	Name   RoomName
}

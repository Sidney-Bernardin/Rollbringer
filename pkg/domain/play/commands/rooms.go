package commands

import (
	"github.com/google/uuid"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/domain/play/results"
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

type CreateRoom struct {
	HostID uuid.UUID
	Name   RoomName
}

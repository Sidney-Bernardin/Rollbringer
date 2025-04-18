package src

import (
	"fmt"

	"github.com/google/uuid"
)

type ExternalErrorType string

const (
	ExternalErrorTypeInternalError ExternalErrorType = "internal_error"
	ExternalErrorTypeUnauthorized  ExternalErrorType = "unauthorized"
)

type ExternalError struct {
	Type    ExternalErrorType `json:"type"`
	Msg     string            `json:"description,omitempty"`
	Details map[string]any    `json:"attrs,omitempty"`
}

func (err *ExternalError) Error() string {
	return fmt.Sprintf("%s: %s", err.Type, err.Msg)
}

/////

type UUID uuid.UUID

const ExternalErrorTypeInvalidUUID ExternalErrorType = "invalid_uuid"

func NewUUID() UUID {
	return UUID(uuid.New())
}

func ParseUUID(str string) (UUID, error) {
	id, err := uuid.Parse(str)
	if err != nil {
		return UUID(uuid.Nil), &ExternalError{
			Type:    ExternalErrorTypeInvalidUUID,
			Msg:     err.Error(),
			Details: map[string]any{"uuid": str},
		}
	}

	return UUID(id), nil
}

func (id UUID) String() string {
	return uuid.UUID(id).String()
}

/////

type RoomUser struct {
	UserID     UUID                `json:"user_id"`
	Permisions []RoomUserPermision `json:"permisions"`
}

type RoomUserPermision string

const (
	RoomUserPermisionOwner      RoomUserPermision = "OWNER"
	RoomUserPermisionGameMaster RoomUserPermision = "GAME_MASTER"
	RoomUserPermisionPlayer     RoomUserPermision = "PLAYER"
)

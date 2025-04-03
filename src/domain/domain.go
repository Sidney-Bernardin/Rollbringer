package domain

import (
	"errors"
	"rollbringer/src"

	"github.com/google/uuid"
)

const (
	ExternalErrorTypeUUIDInvalid src.ExternalErrorType = "uuid_invalid"
	ExternalErrorTypeViewInvalid src.ExternalErrorType = "view_invalid"
)

var (
	ErrEntityConflict     = errors.New("entity conflict")
	ErrEntityNotFound     = errors.New("entity not found")
	ErrNoEntitiesEffected = errors.New("no entities effected")
)

/////

type UUID uuid.UUID

func NewUUID() UUID {
	return UUID(uuid.New())
}

func ParseUUID(str string) (UUID, error) {
	id, err := uuid.Parse(str)
	if err != nil {
		return UUID(uuid.Nil), &src.ExternalError{
			Type:        ExternalErrorTypeUUIDInvalid,
			Description: err.Error(),
			Details:     map[string]any{"uuid": str},
		}
	}

	return UUID(id), nil
}

func (id UUID) String() string {
	return uuid.UUID(id).String()
}

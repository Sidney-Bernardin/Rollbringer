package domain

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrEntityConflict     = errors.New("entity conflict")
	ErrEntityNotFound     = errors.New("entity not found")
	ErrNoEntitiesEffected = errors.New("no entities effected")
	ErrInvalidView        = errors.New("invalid view")
)

/////

type UUID uuid.UUID

func ParseUUID(str string) (UUID, error) {
	id, err := uuid.Parse(str)
	if err != nil {
		return UUID(uuid.Nil), &DomainError{
			Type:        DomainErrorTypeUUIDInvalid,
			Description: err.Error(),
			Attrs:       map[string]any{"uuid": str},
		}
	}

	return UUID(id), nil
}

func (id UUID) String() string {
	return uuid.UUID(id).String()
}

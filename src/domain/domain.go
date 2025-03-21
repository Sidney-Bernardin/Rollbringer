package domain

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrEntityConflict     = errors.New("entity conflict")
	ErrEntityNotFound     = errors.New("entity not found")
	ErrNoEntitiesEffected = errors.New("no entities effected")
	ErrInvalidView        = errors.New("invalid view")
)

/////

type DomainErrorType string

const (
	DomainErrorTypeUUIDInvalid    = "uuid_invalid"
	DomainErrorTypeEntityNotFound = "entity_not_found"
)

type DomainError struct {
	Type        DomainErrorType
	Description string
	Details     map[string]any
}

func (err *DomainError) Error() string {
	return fmt.Sprintf("%s: %s", err.Type, err.Description)
}

/////

type UUID uuid.UUID

func ParseUUID(str string) (UUID, error) {
	id, err := uuid.Parse(str)
	if err != nil {
		return UUID(uuid.Nil), &DomainError{
			Type:        DomainErrorTypeUUIDInvalid,
			Description: err.Error(),
			Details:     map[string]any{"uuid": str},
		}
	}

	return UUID(id), nil
}

func (id UUID) String() string {
	return uuid.UUID(id).String()
}

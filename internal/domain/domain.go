package domain

import (
	"fmt"
	"reflect"

	"github.com/google/uuid"
	googleUUID "github.com/google/uuid"
)

type UserError struct {
	Type    UserErrorType  `json:"type"`
	Message string         `json:"message,omitempty"`
	Details map[string]any `json:"details,omitempty"`
}

func (err *UserError) Error() string {
	return fmt.Sprintf("%s: %s %v", err.Type, err.Message, err.Details)
}

type UserErrorType string

const (
	DomainErrorTypeInternalServerError UserErrorType = "internal-server-error"
	DomainErrorTypeUUIDInvalid         UserErrorType = "uuid-invalid"
)

/////

type EntityConflictError[E any] struct {
	Column  string
	Message string
}

func (err *EntityConflictError[E]) Error() string {
	return fmt.Sprintf("%s conflict on %s: %v", reflect.TypeFor[E]().Name(), err.Column, err.Message)
}

type NoEntitiesError[E any] struct{}

func (err *NoEntitiesError[E]) Error() string {
	return fmt.Sprintf("no %s rows", reflect.TypeFor[E]().Name())
}

type NoEntitiesEffectedError[E any] struct{}

func (err *NoEntitiesEffectedError[E]) Error() string {
	return fmt.Sprintf("no %s rows effected", reflect.TypeFor[E]().Name())
}

/////

type UUID struct {
	googleUUID.UUID
}

func NewRandomUUID() UUID {
	return UUID{uuid.Must(uuid.NewRandom())}
}

func ParseUUID(uuid string) (ret UUID, err error) {
	gUUID, err := googleUUID.Parse(uuid)
	if err != nil {
		return ret, &UserError{
			Type:    DomainErrorTypeUUIDInvalid,
			Message: err.Error(),
		}
	}
	return UUID{gUUID}, nil
}

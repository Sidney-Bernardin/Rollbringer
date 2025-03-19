package domain

import (
	"errors"
	"fmt"
)

var (
	ErrEntityConflict     = errors.New("entity conflict")
	ErrEntityNotFound     = errors.New("entity not found")
	ErrNoEntitiesEffected = errors.New("no entities effected")
	ErrInvalidView        = errors.New("invalid view")
)

type DomainErrorType string

type DomainError struct {
	Type        DomainErrorType
	Description string
	Details     map[string]any
}

func (err *DomainError) Error() string {
	return fmt.Sprintf("%s: %s", err.Type, err.Description)
}

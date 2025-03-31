package domain

import (
	"errors"
)

var (
	ErrEntityConflict     = errors.New("entity conflict")
	ErrEntityNotFound     = errors.New("entity not found")
	ErrNoEntitiesEffected = errors.New("no entities effected")
	ErrInvalidView        = errors.New("invalid view")
)

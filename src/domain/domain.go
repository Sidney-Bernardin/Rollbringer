package domain

import (
	"errors"
	"rollbringer/src"
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

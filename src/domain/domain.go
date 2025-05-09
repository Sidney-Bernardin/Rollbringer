package domain

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrEntityConflict     = errors.New("entity conflict")
	ErrEntityNotFound     = errors.New("entity not found")
	ErrNoEntitiesEffected = errors.New("no entities effected")
)

type PublicBroker interface {
	Pub(ctx context.Context, event any) bool
	SubUser(ctx context.Context, userID uuid.UUID, callback func(event any)) error
	SubRoom(ctx context.Context, roomID uuid.UUID, callback func(event any)) error
}

type PublicUser struct {
	UserID         uuid.UUID `json:"user_id"`
	Username       string    `json:"username"`
	ProfilePicture string    `json:"profile_picture"`
}

type ExternalError struct {
	Type    ExternalErrorType `json:"type"`
	Msg     string            `json:"description,omitempty"`
	Details map[string]any    `json:"attrs,omitempty"`
}

type ExternalErrorType string

const (
	ExternalErrorTypeInternalError ExternalErrorType = "internal-error"
	ExternalErrorTypeUnauthorized  ExternalErrorType = "unauthorized"
	ExternalErrorTypeInvalidUUID   ExternalErrorType = "invalid-uuid"
)

func (err *ExternalError) Error() string {
	return fmt.Sprintf("%s: %s", err.Type, err.Msg)
}

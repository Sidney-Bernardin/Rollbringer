package domain

import (
	"context"
	"errors"

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

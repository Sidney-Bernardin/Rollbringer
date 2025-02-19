package service

import (
	"context"
	"rollbringer/pkg/domain"

	"github.com/google/uuid"
)

type PlayDatabaseRepository interface {
	Close() error
	Transaction(context.Context, func(PlayDatabaseRepository) error) error

	RoomInsert(ctx context.Context, room *domain.Room) error
	RoomGet(ctx context.Context, key string, value any) (*domain.Room, error)
	RoomsGet(ctx context.Context, key string, value any) ([]*domain.Room, error)
	RoomDelete(ctx context.Context, roomID uuid.UUID, ownerID uuid.UUID) error

	BoardInsert(ctx context.Context, view domain.BoardView, board *domain.Board) error
	BoardGet(ctx context.Context, view domain.BoardView, key string, value any) (*domain.Board, error)
	BoardsGet(ctx context.Context, view domain.BoardView, key string, value any) ([]*domain.Board, error)
}

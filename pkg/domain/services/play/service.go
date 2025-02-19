package service

import (
	"context"
	"log/slog"

	"rollbringer/pkg/domain"

	"github.com/google/uuid"
)

type PlayService interface {
	domain.IService

	CreateRoom(ctx context.Context, session *domain.Session, room *domain.Room) error
	GetRoom(ctx context.Context, roomID uuid.UUID) (*domain.Room, error)
	GetRooms(ctx context.Context, ownerID uuid.UUID) ([]*domain.Room, error)
	DeleteRoom(ctx context.Context, session *domain.Session, roomID uuid.UUID) error

	CreateBoard(ctx context.Context, session *domain.Session, view domain.BoardView, board *domain.Board) error
	GetBoard(ctx context.Context, view domain.BoardView, boardID uuid.UUID) (*domain.Board, error)
}

type playService struct {
	*domain.Service

	playDBRepo PlayDatabaseRepository
}

func New(
	config *domain.Config,
	logger *slog.Logger,
	pubSub domain.PubSubRepository,
	playDBRepo PlayDatabaseRepository,
) PlayService {
	return &playService{
		Service: &domain.Service{
			Config: config,
			Logger: logger,
			PubSub: pubSub,
		},
		playDBRepo: playDBRepo,
	}
}

func (svc *playService) Run(ctx context.Context) error {
	err := svc.PubSub.Subscribe(ctx, "play", svc.subPlay, map[domain.Operation]any{
		domain.OperationGetRoomRequest:  &domain.GetRoomRequest{},
		domain.OperationGetRoomsRequest: &domain.GetRoomsRequest{},
	})

	return domain.Wrap(err, "cannot subscribe", map[string]any{"subject": "play"})
}

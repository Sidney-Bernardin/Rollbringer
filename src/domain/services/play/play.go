package play

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"rollbringer/src"
	"rollbringer/src/domain"
)

type Service interface {
	CreateRoom(ctx context.Context, creatorID uuid.UUID, args *CreateRoomOpts) (*Room, error)
	CreateBoard(ctx context.Context, opts *CreateBoardOpts, creator *domain.PublicUser, users []domain.PublicUser) (*Board, error)
	JoinRoom(ctx context.Context, roomID uuid.UUID, newcomer *domain.PublicUser) (*Room, error)
}

type (
	Database interface {
		DatabaseCommon
		CreateRoom(ctx context.Context, room *Room) error
		CreateBoard(ctx context.Context, board *Board) error
		JoinRoom(ctx context.Context, userID, roomID uuid.UUID, permisions ...RoomUserPermision) (*Room, bool, error)
	}

	DatabaseCommon interface {
		GetRoomsByUserID(ctx context.Context, userID uuid.UUID) ([]*Room, error)
		GetUserBoard(ctx context.Context, userID, boardID uuid.UUID) (*Board, error)
		GetBoardsByUserID(ctx context.Context, boardID uuid.UUID) ([]*Board, error)
	}
)

type service struct {
	config *src.Config
	log    *slog.Logger

	broker domain.Broker
	db     Database
}

func NewService(config *src.Config, log *slog.Logger, broker domain.Broker, database Database) Service {
	return &service{config, log, broker, database}
}

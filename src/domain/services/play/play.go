package play

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"rollbringer/src"

	"rollbringer/src/domain"
)

type EventChat struct {
	RoomID   string `json:"room_id"`
	AuthorID string `json:"author_id"`
	Message  string `json:"message"`
}

type Service interface {
	CreateRoom(ctx context.Context, creatorID uuid.UUID, args *CreateRoomOpts) (*Room, error)
	CreateBoard(ctx context.Context, opts *CreateBoardOpts, creator *domain.EventNewBoardUser, users []domain.EventNewBoardUser) (*Board, error)
	JoinRoom(ctx context.Context, roomID uuid.UUID, newcomer *domain.EventRoomJoinedNewcomer) (*Room, error)
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
		GetBoardsByUserID(ctx context.Context, boardID uuid.UUID) ([]*Board, error)
	}
)

type (
	Broker interface {
		BrokerCommon
	}

	BrokerCommon interface {
		domain.PublicBroker
		SubChat(ctx context.Context, roomID uuid.UUID, callback func(event *EventChat)) error
	}
)

type service struct {
	config *src.Config
	log    *slog.Logger

	broker   Broker
	database Database
}

func NewService(config *src.Config, log *slog.Logger, broker Broker, database Database) Service {
	return &service{config, log, broker, database}
}

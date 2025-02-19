package service

import (
	"context"

	"github.com/google/uuid"

	"rollbringer/pkg/domain"
)

func (svc *playService) CreateBoard(ctx context.Context, session *domain.Session, view domain.BoardView, board *domain.Board) error {
	err := svc.playDBRepo.BoardInsert(ctx, view, board)
	return domain.Wrap(err, "cannot insert board", nil)
}

func (svc *playService) GetBoard(ctx context.Context, view domain.BoardView, boardID uuid.UUID) (*domain.Board, error) {
	board, err := svc.playDBRepo.BoardGet(ctx, view, "id", boardID)
	return board, domain.Wrap(err, "cannot get board", nil)
}

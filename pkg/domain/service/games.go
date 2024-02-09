package service

import (
	"context"

	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
)

func (svc *Service) CreateGame(ctx context.Context, userID string) (string, string, error) {
	gameID, title, err := svc.db.InsertGame(ctx, userID)
	return gameID, title, errors.Wrap(err, "cannot insert game")
}

func (svc *Service) GetGame(ctx context.Context, gameID string) (*domain.Game, error) {
	game, err := svc.db.GetGame(ctx, gameID)
	return game, errors.Wrap(err, "cannot get game")
}

func (svc *Service) GetGamesFromUser(ctx context.Context, userID string) ([]*domain.Game, error) {
	games, err := svc.db.GetGamesFromUser(ctx, userID)
	return games, errors.Wrap(err, "cannot get games from user")
}

func (svc *Service) DeleteGame(ctx context.Context, gameID, userID string) error {
	err := svc.db.DeleteGame(ctx, gameID, userID)
	return errors.Wrap(err, "cannot delete game")
}

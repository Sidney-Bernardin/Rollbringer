package service

import (
	"context"

	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
)

func (svc *Service) CreateGame(ctx context.Context, session *domain.Session) (*domain.Game, error) {
	domain.ParseUUIDs(&session.ID, &session.UserID)

	// Create a game.
	game := &domain.Game{
		HostID: session.UserID,
		Title:  "New Game %d",
	}

	// Insert the game.
	if err := svc.DB.InsertGame(ctx, game); err != nil {
		return nil, errors.Wrap(err, "cannot insert game")
	}

	return game, nil
}

func (svc *Service) GetGames(ctx context.Context, userID string) ([]*domain.Game, error) {
	domain.ParseUUIDs(&userID)
	games, err := svc.DB.GetGames(ctx, userID)
	return games, errors.Wrap(err, "cannot get games from user")
}

func (svc *Service) GetGame(ctx context.Context, gameID string) (*domain.Game, error) {
	domain.ParseUUIDs(&gameID)
	game, err := svc.DB.GetGame(ctx, gameID)
	return game, errors.Wrap(err, "cannot get game")
}

func (svc *Service) DeleteGame(ctx context.Context, gameID, userID string) error {
	domain.ParseUUIDs(&gameID, &userID)
	err := svc.DB.DeleteGame(ctx, gameID, userID)
	return errors.Wrap(err, "cannot delete game")
}

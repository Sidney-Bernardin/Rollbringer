package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
)

func (svc *Service) CreateGame(ctx context.Context, userID string) (string, string, error) {

	// Parse the user-ID.
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return "", "", domain.ErrUserNotFound
	}

	// Insert a new game.
	gameID, title, err := svc.db.InsertGame(ctx, parsedUserID)
	return gameID.String(), title, errors.Wrap(err, "cannot insert game")
}

func (svc *Service) GetGame(ctx context.Context, gameID string) (*domain.Game, error) {

	// Parse the game-ID.
	parsedGameID, err := uuid.Parse(gameID)
	if err != nil {
		return nil, domain.ErrGameNotFound
	}

	// Get the game.
	game, err := svc.db.GetGame(ctx, parsedGameID)
	return game, errors.Wrap(err, "cannot get game")
}

func (svc *Service) GetGamesFromUser(ctx context.Context, userID string) ([]*domain.Game, error) {

	// Parse the user-ID.
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// Get the games.
	games, err := svc.db.GetGamesFromUser(ctx, parsedUserID)
	return games, errors.Wrap(err, "cannot get games from user")
}

func (svc *Service) DeleteGame(ctx context.Context, gameID, userID string) error {

	// Parse the game-ID.
	parsedGameID, err := uuid.Parse(gameID)
	if err != nil {
		return domain.ErrGameNotFound
	}

	// Parse the user-ID.
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	// Delete the game.
	err = svc.db.DeleteGame(ctx, parsedGameID, parsedUserID)
	return errors.Wrap(err, "cannot delete game")
}

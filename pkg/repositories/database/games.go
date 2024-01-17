package database

import (
	"context"
	"database/sql"
	"fmt"
	"rollbringer/pkg/models"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (database *Database) CreateGame(ctx context.Context, userID uuid.UUID) (uuid.UUID, string, error) {

	// Get the count of the user's games.
	var count int
	q := `SELECT COUNT(*) FROM games WHERE host_id=$1`
	if err := database.db.QueryRowContext(ctx, q, userID).Scan(&count); err != nil {
		return uuid.Nil, "", errors.Wrap(err, "cannot select games")
	}

	// Verify the amount of games the user has.
	if count >= 5 {
		return uuid.Nil, "", ErrMaxGames
	}

	// Insert the a new game into the database.
	gameID := uuid.New()
	title := fmt.Sprintf("Untitled %d", count+1)
	q = `INSERT INTO games (id, host_id, title) VALUES ($1, $2, $3)`
	if _, err := database.db.ExecContext(ctx, q, gameID, userID, title); err != nil {
		return uuid.Nil, "", errors.Wrap(err, "cannot insert game")
	}

	return gameID, title, nil
}

func (database *Database) GetGame(ctx context.Context, gameID uuid.UUID) (*models.Game, error) {

	var game models.Game
	q := `SELECT id, host_id, title FROM games WHERE id=$1`
	err := database.db.QueryRowContext(ctx, q, gameID).Scan(&game.ID, &game.HostID, &game.Title)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, ErrGameNotFound
		}

		return nil, errors.Wrap(err, "cannot select user")
	}

	return &game, nil
}

func (database *Database) GetGames(ctx context.Context, hostID uuid.UUID) ([]*models.Game, error) {
	games := []*models.Game{}

	// Get the host's games from the database.
	q := `SELECT id, host_id, title FROM games WHERE host_id=$1`
	rows, err := database.db.QueryContext(ctx, q, hostID)
	if err != nil {

		if err == sql.ErrNoRows {
			return games, nil
		}

		return nil, errors.Wrap(err, "cannot select games")
	}

	for rows.Next() {

		// Scan the current row and append it to the games slice.
		var game models.Game
		if err := rows.Scan(&game.ID, &game.HostID, &game.Title); err != nil {
			return nil, errors.Wrap(err, "cannot scan game")
		}
		games = append(games, &game)
	}

	return games, nil
}

func (database *Database) DeleteGame(ctx context.Context, userID uuid.UUID, gameID uuid.UUID) error {

	q := `DELETE FROM games WHERE id=$1`
	if _, err := database.db.ExecContext(ctx, q, gameID); err != nil {

		if err == sql.ErrNoRows {
			return ErrUnauthorized
		}

		return errors.Wrap(err, "cannot delete game")
	}

	return nil
}

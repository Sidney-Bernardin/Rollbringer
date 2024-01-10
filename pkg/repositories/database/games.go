package database

import (
	"context"
	"database/sql"
	"rollbringer/pkg/models"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (database *Database) CreateGame(ctx context.Context, session *models.Session) (string, error) {
	return "", nil
}

func (database *Database) GetGame(ctx context.Context, gameID uuid.UUID) (*models.Game, error) {
	var game models.Game

	q := `SELECT id, host_id, title, player_ids FROM games WHERE id=$1`
	err := database.db.QueryRowContext(ctx, q, gameID).Scan(&game.ID, &game.HostID, &game.Title, &game.PlayerIDs)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrGameNotFound
		}

		return nil, errors.Wrap(err, "cannot select user")
	}

	return &game, nil
}

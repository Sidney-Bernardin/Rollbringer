package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
)

// InsertGame inserts a new game for the host. If the host has more than 5
// games, returns domain.ErrMaxGames
func (db *Database) InsertGame(ctx context.Context, hostID string) (string, string, error) {

	hostUUID, _ := uuid.Parse(hostID)

	// Get the number of games with the host-ID.
	var count int
	err := db.conn.
		QueryRow(ctx, `SELECT COUNT(*) FROM games WHERE host_id = $1`, hostUUID).
		Scan(&count)

	if err != nil {
		return "", "", errors.Wrap(err, "cannot count games")
	}

	if count >= 5 {
		return "", "", domain.ErrMaxGames
	}

	gameID := uuid.New().String()
	title := fmt.Sprintf("Game %d", count+1)

	// Insert a new game for the host.
	_, err = db.conn.Exec(ctx,
		`INSERT INTO games (id, host_id, title) VALUES ($1, $2, $3)`,
		gameID, hostUUID, title)

	if err != nil {
		return "", "", errors.Wrap(err, "cannot insert game")
	}

	return gameID, title, nil
}

// GetGame returns the game with the game-ID from the database. If the game
// doesn't exist, returns domain.ErrGameNotFound.
func (db *Database) GetGame(ctx context.Context, gameID string) (*domain.Game, error) {

	gameUUID, _ := uuid.Parse(gameID)

	// Get the game with the game-ID.
	rows, err := db.conn.Query(ctx, `SELECT * FROM games WHERE id = $1`, gameUUID)
	if err != nil {
		return nil, errors.Wrap(err, "cannot select game")
	}
	defer rows.Close()

	// Scan into a game model.
	game, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[domain.Game])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrGameNotFound
		}

		return nil, errors.Wrap(err, "cannot scan game")
	}

	return game, nil
}

// GetGames return the games with the host-ID from the database.
func (db *Database) GetGamesFromUser(ctx context.Context, hostID string) ([]*domain.Game, error) {

	hostUUID, _ := uuid.Parse(hostID)

	// Get the games with the host-ID.
	rows, err := db.conn.Query(ctx, `SELECT * FROM games WHERE host_id = $1`, hostUUID)
	if err != nil {
		return nil, errors.Wrap(err, "cannot select games")
	}
	defer rows.Close()

	// Scan into a slice of game models.
	games, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[domain.Game])
	return games, errors.Wrap(err, "cannot scan games")
}

// DeleteGame deletes the game with the game-ID and host-ID from the database.
// If the game doesn't exist, returns domain.ErrGameNotFound.
func (db *Database) DeleteGame(ctx context.Context, gameID, hostID string) error {

	gameUUID, _ := uuid.Parse(gameID)
	hostUUID, _ := uuid.Parse(hostID)

	// Delete the game with the game-ID and host-ID.
	cmdTag, err := db.conn.Exec(ctx,
		`DELETE FROM games WHERE id = $1 AND host_id = $2`,
		gameUUID, hostUUID)

	if err != nil {
		return errors.Wrap(err, "cannot delete game")
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.ErrGameNotFound
	}

	return nil
}

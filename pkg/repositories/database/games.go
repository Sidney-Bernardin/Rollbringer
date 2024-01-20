package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	"rollbringer/pkg/models"
)

// InsertGame inserts a new game for the host. If the host has more than 5
// games, returns ErrMaxGames
func (db *Database) InsertGame(ctx context.Context, hostID uuid.UUID) (uuid.UUID, string, error) {

	// Get the number of games with the host-ID.
	var count int
	err := db.conn.
		QueryRow(ctx, `SELECT COUNT(*) FROM games WHERE host_id = $1`, hostID).
		Scan(&count)

	if err != nil {
		return uuid.Nil, "", errors.Wrap(err, "cannot count games")
	}

	if count >= 5 {
		return uuid.Nil, "", ErrMaxGames
	}

	gameID := uuid.New()
	title := fmt.Sprintf("Game %d", count+1)

	// Insert a new game for the host.
	_, err = db.conn.Exec(ctx,
		`INSERT INTO games (id, host_id, title) VALUES ($1, $2, $3)`,
		gameID, hostID, title)

	if err != nil {
		return uuid.Nil, "", errors.Wrap(err, "cannot insert game")
	}

	return gameID, title, nil
}

// GetGame returns the game with the game-ID from the database. If the game
// doesn't exist, returns ErrGameNotFound.
func (db *Database) GetGame(ctx context.Context, gameID uuid.UUID) (*models.Game, error) {

	// Get the game with the game-ID.
	rows, err := db.conn.Query(ctx, `SELECT * FROM games WHERE id = $1`, gameID)
	if err != nil {
		return nil, errors.Wrap(err, "cannot select game")
	}
	defer rows.Close()

	// Scan into a game model.
	game, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[models.Game])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrGameNotFound
		}

		return nil, errors.Wrap(err, "cannot scan game")
	}

	return game, nil
}

// GetGames return the games with the host-ID from the database.
func (db *Database) GetGames(ctx context.Context, hostID uuid.UUID) ([]*models.Game, error) {

	// Get the games with the host-ID.
	rows, err := db.conn.Query(ctx, `SELECT * FROM games WHERE host_id = $1`, hostID)
	if err != nil {
		return []*models.Game{}, errors.Wrap(err, "cannot select games")
	}
	defer rows.Close()

	// Scan into a slice of game models.
	games, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[models.Game])
	if err != nil {
		return []*models.Game{}, errors.Wrap(err, "cannot scan games")
	}

	return games, nil
}

// DeleteGame deletes the game with the game-ID and host-ID from the database.
// If the game doesn't exist, returns ErrGameNotFound.
func (db *Database) DeleteGame(ctx context.Context, gameID uuid.UUID, hostID uuid.UUID) error {

	// Delete the game with the game-ID and host-ID.
	cmdTag, err := db.conn.Exec(ctx,
		`DELETE FROM games WHERE id = $1 AND host_id = $2`,
		gameID, hostID)

	if err != nil {
		return errors.Wrap(err, "cannot delete game")
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrGameNotFound
	}

	return nil
}

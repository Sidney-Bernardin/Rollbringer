package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
)

// InsertGame inserts the game. If the host has more than 5 games,
// returns domain.ErrMaxGames
func (db *Database) InsertGame(ctx context.Context, game *domain.Game) error {
	db.parseUUIDs(&game.HostID)

	// Get the number of games with the host-ID.
	var count int
	err := db.conn.QueryRow(
		`SELECT COUNT(*) FROM games WHERE host_id = $1`, game.HostID).
		Scan(&count)

	if err != nil {
		return errors.Wrap(err, "cannot select games count")
	}

	if count >= 5 {
		return domain.ErrMaxGames
	}

	game.ID = uuid.New().String()
	game.Title = fmt.Sprintf(game.Title, count+1)
	game.PDFs = []string{}

	// Insert the game.
	_, err = db.conn.Exec(
		`INSERT INTO games (id, host_id, title, pdfs) VALUES ($1, $2, $3, $4)`,
		game.ID, game.HostID, game.Title, pq.Array(game.PDFs))

	return errors.Wrap(err, "cannot insert game")
}

// GetGames return the games with the host-ID.
func (db *Database) GetGames(ctx context.Context, hostID string) ([]*domain.Game, error) {
	db.parseUUIDs(&hostID)

	// Get the games with the host-ID.
	rows, err := db.conn.Query(
		`SELECT id, host_id, title, pdfs FROM games WHERE host_id = $1`, hostID)

	if err != nil {
		return nil, errors.Wrap(err, "cannot select games")
	}
	defer rows.Close()

	// Scan the rows into a slice of games.
	games := []*domain.Game{}
	for rows.Next() {
		var game domain.Game
		if err := rows.Scan(&game.ID, &game.HostID, &game.Title, pq.Array(&game.PDFs)); err != nil {
			return nil, errors.Wrap(err, "cannot scan game")
		}
		games = append(games, &game)
	}

	return games, nil
}

// GetGame returns the game with the game-ID. If the game doesn't exist,
// returns domain.ErrGameNotFound.
func (db *Database) GetGame(ctx context.Context, gameID string) (*domain.Game, error) {
	db.parseUUIDs(&gameID)

	// Get the game with the game-ID.
	var game domain.Game
	err := db.conn.QueryRow(
		`SELECT id, host_id, title, pdfs FROM games WHERE id = $1`, gameID).
		Scan(&game.ID, &game.HostID, &game.Title, pq.Array(&game.PDFs))

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrGameNotFound
		}

		return nil, errors.Wrap(err, "cannot select game")
	}

	return &game, nil
}

// AppendGamePDF appends the PDF-ID to the game with the game-ID.
func (db *Database) AppendGamePDF(ctx context.Context, gameID, pdfID string) error {
	db.parseUUIDs(&gameID, &pdfID)

	// Append the PDF to the game.
	result, err := db.conn.Exec(
		`UPDATE games SET pdfs = array_append(pdfs, $1) WHERE id = $2`,
		pdfID, gameID)

	if err != nil {
		return errors.Wrap(err, "cannot append pdf to game")
	}

	// Get the number of rows affected.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "cannot get rows affected")
	}

	if rowsAffected == 0 {
		return domain.ErrGameNotFound
	}

	return nil
}

func (db *Database) RemoveGamePDF(ctx context.Context, gameID, pdfID string) error {
	db.parseUUIDs(&gameID, &pdfID)

	// Remove the PDF from the game.
	result, err := db.conn.Exec(
		`UPDATE games SET pdfs = array_remove(pdfs, $1) WHERE id = $2`,
		pdfID, gameID)

	if err != nil {
		return errors.Wrap(err, "cannot remove pdf from game")
	}

	// Get the number of rows affected.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "cannot get rows affected")
	}

	if rowsAffected == 0 {
		return domain.ErrGameNotFound
	}

	return nil
}

// DeleteGame deletes the game with the game-ID and host-ID. If the game doesn't
// exist, returns domain.ErrGameNotFound.
func (db *Database) DeleteGame(ctx context.Context, gameID, hostID string) error {
	db.parseUUIDs(&gameID, &hostID)

	// Delete the game with the game-ID and host-ID.
	result, err := db.conn.Exec(
		`DELETE FROM games WHERE id = $1 AND host_id = $2`,
		gameID, hostID)

	if err != nil {
		return errors.Wrap(err, "cannot delete game")
	}

	// Get the number of rows affected.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "cannot get rows affected")
	}

	if rowsAffected == 0 {
		return domain.ErrGameNotFound
	}

	return nil
}

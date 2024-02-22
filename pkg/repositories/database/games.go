package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
)

// InsertGame inserts the game. If the host has more than 5 games,
// returns domain.ErrMaxGames
func (db *Database) InsertGame(ctx context.Context, game *domain.Game) error {

	hostUUID, _ := uuid.Parse(game.HostID)

	// Get the number of games with the host-ID.
	rows, err := db.conn.Query(ctx, `SELECT COUNT(*) FROM games WHERE host_id = $1`, hostUUID)
	if err != nil {
		return errors.Wrap(err, "cannot select games")
	}
	defer rows.Close()

	count, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return errors.Wrap(err, "cannot scan count")
	}

	if count >= 5 {
		return domain.ErrMaxGames
	}

	game.ID = uuid.New().String()
	game.Title = fmt.Sprintf(game.Title, count+1)

	// Insert the game.
	_, err = db.conn.Exec(ctx,
		`INSERT INTO games (id, host_id, title) VALUES ($1, $2, $3)`,
		game.ID, hostUUID, game.Title)

	return errors.Wrap(err, "cannot insert game")
}

// InsertRoll inserts the roll.
func (db *Database) InsertRoll(ctx context.Context, roll *domain.Roll) error {

	gameUUID, _ := uuid.Parse(roll.GameID)
	roll.ID = uuid.New().String()

	// Insert the roll.
	_, err := db.conn.Exec(ctx,
		`INSERT INTO rolls (id, game_id, die_expressions, die_results, modifier_expression, modifier_result)
			VALUES ($1, $2, $3, $4, $5, $6)`,
		roll.ID, gameUUID, roll.DieExpressions, roll.DieResults, roll.ModifierExpression, roll.ModifierResult,
	)

	return errors.Wrap(err, "cannot insert roll")
}

// GetGames return the games with the host-ID.
func (db *Database) GetGames(ctx context.Context, hostID string) ([]*domain.Game, error) {

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

// GetGame returns the game with the game-ID. If the game doesn't exist,
// returns domain.ErrGameNotFound.
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

// GetRolls returns the rolls with the game-ID.
func (db *Database) GetRolls(ctx context.Context, gameID string) ([]*domain.Roll, error) {

	gameUUID, _ := uuid.Parse(gameID)

	// Get the rolls with the game-ID.
	rows, err := db.conn.Query(ctx, `SELECT * FROM rolls WHERE game_id = $1`, gameUUID)
	if err != nil {
		return nil, errors.Wrap(err, "cannot select rolls")
	}
	defer rows.Close()

	// Scan into a slice of roll models.
	rolls, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[domain.Roll])
	return rolls, errors.Wrap(err, "cannot scan rolls")
}

// DeleteGame deletes the game with the game-ID and host-ID. If the game doesn't
// exist, returns domain.ErrGameNotFound.
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

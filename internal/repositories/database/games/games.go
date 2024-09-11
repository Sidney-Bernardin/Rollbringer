package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"rollbringer/internal"
	"rollbringer/internal/repositories/databases"
)

func gameColumns(views map[string]internal.GameView) (columns string) {
	switch views["game"] {
	case internal.GameViewGameAll:
		columns += `games.*`
	}

	switch views["host"] {
	case internal.GameViewHostInfo:
		columns += `, users.id AS "host.id", users.username AS "host.username"`
	}

	return columns
}

func gameJoins(views map[string]internal.GameView) (joins string) {
	if _, ok := views["host"]; ok {
		joins += `LEFT JOIN users ON users.id = games.host_id`
	}

	return joins
}

func (db *gamesSchema) GameInsert(ctx context.Context, game *internal.Game) error {
	err := db.TX.QueryRowxContext(ctx,
		`INSERT INTO games (id, host_id, name)
			VALUES ($1, $2, $3)
		RETURNING id`,
		uuid.New(), game.HostID, game.Name,
	).Scan(&game.ID)

	return errors.Wrap(err, "cannot insert game")
}

func (db *gamesSchema) GamesCount(ctx context.Context, hostID uuid.UUID) (count int, err error) {
	err = sqlx.GetContext(ctx, db.TX, &count,
		`SELECT COUNT(*) FROM games WHERE host_id = $1`, hostID)

	return count, errors.Wrap(err, "cannot count games")
}

func (db *gamesSchema) GameGet(ctx context.Context, gameID uuid.UUID, views map[string]internal.GameView) (*internal.Game, error) {

	var game databases.Game
	query := fmt.Sprintf(`SELECT %s FROM games %s WHERE games.id = $1`, gameColumns(views), gameJoins(views))
	if err := sqlx.GetContext(ctx, db.TX, &game, query, gameID); err != nil {
		if err == sql.ErrNoRows {
			return nil, internal.NewProblemDetail(ctx, internal.PDOpts{
				Type:   internal.PDTypeGameNotFound,
				Detail: "Can't find a game with the given game_id.",
				Extra: map[string]any{
					"game_id": gameID,
				},
			})
		}

		return nil, errors.Wrap(err, "cannot select game")
	}

	return game.Internalized(), nil
}

func (db *gamesSchema) GamesGetForHost(ctx context.Context, hostID uuid.UUID, views map[string]internal.GameView) ([]*internal.Game, error) {

	var games []*databases.Game
	query := fmt.Sprintf(`SELECT %s FROM games %s WHERE games.host_id = $1`, gameColumns(views), gameJoins(views))
	if err := sqlx.SelectContext(ctx, db.TX, &games, query, hostID); err != nil {
		return nil, errors.Wrap(err, "cannot select games")
	}

	// Internalize each game.
	ret := make([]*internal.Game, len(games))
	for i, m := range games {
		ret[i] = m.Internalized()
	}

	return ret, nil
}

func (db *gamesSchema) GameDelete(ctx context.Context, gameID, hostID uuid.UUID) error {
	_, err := db.TX.ExecContext(ctx,
		`DELETE FROM games WHERE id = $1 AND host_id = $2`, gameID, hostID)

	return errors.Wrap(err, "cannot delete game")
}

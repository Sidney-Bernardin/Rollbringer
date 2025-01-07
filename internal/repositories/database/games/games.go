package games

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"rollbringer/internal"
	"rollbringer/internal/repositories/database"
)

func gameColumns(view internal.GameView) string {
	switch view {
	case internal.GameViewListItem:
		return `games.*,` +
			`users.id AS "host.id",` +
			`users.username AS "host.username"`

	default:
		return `games.*`
	}
}

func gameJoins(view internal.GameView) string {
	switch view {
	case internal.GameViewListItem:
		return `LEFT JOIN users.users ON users.id = games.host_id`
	default:
		return ``
	}
}

func (db *gamesSchema) GameInsert(ctx context.Context, game *internal.Game) error {
	err := db.TX.QueryRowxContext(ctx,
		`INSERT INTO games.games (id, host_id, name)
			VALUES ($1, $2, $3)
		RETURNING id`,
		uuid.New(), game.HostID, game.Name,
	).Scan(&game.ID)

	return errors.Wrap(err, "cannot insert game")
}

func (db *gamesSchema) GamesCount(ctx context.Context, hostID uuid.UUID) (count int, err error) {
	err = sqlx.GetContext(ctx, db.TX, &count,
		`SELECT COUNT(*) FROM games.games WHERE host_id = $1`, hostID)

	return count, errors.Wrap(err, "cannot count games")
}

func (db *gamesSchema) GameGet(ctx context.Context, gameID uuid.UUID, view internal.GameView) (*internal.Game, error) {

	var game database.Game
	query := fmt.Sprintf(`SELECT %s FROM games.games %s WHERE games.id = $1`, gameColumns(view), gameJoins(view))
	if err := sqlx.GetContext(ctx, db.TX, &game, query, gameID); err != nil {
		if err == sql.ErrNoRows {
			return nil, internal.NewProblemDetail(ctx, internal.PDOpts{
				Type:   internal.PDTypeGameNotFound,
				Detail: "Cannot find a game with the given game_id.",
				Extra: map[string]any{
					"game_id": gameID,
				},
			})
		}

		return nil, errors.Wrap(err, "cannot select game")
	}

	return game.Internalized(), nil
}

func (db *gamesSchema) GamesGetByHost(ctx context.Context, hostID uuid.UUID, view internal.GameView) ([]*internal.Game, error) {

	var games []*database.Game
	query := fmt.Sprintf(`SELECT %s FROM games.games %s WHERE games.host_id = $1`, gameColumns(view), gameJoins(view))
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

func (db *gamesSchema) GamesGetByUser(ctx context.Context, userID uuid.UUID, view internal.GameView) ([]*internal.Game, error) {
	query := fmt.Sprintf(` 
		SELECT %s FROM games.games %s
		WHERE EXISTS (
			SELECT * FROM game_users WHERE game_users.user_id = $1 AND game_users.game_id = games.id
		)
	`, gameColumns(view), gameJoins(view))

	var games []*database.Game
	if err := sqlx.SelectContext(ctx, db.TX, &games, query, userID); err != nil {
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
		`DELETE FROM games.games WHERE id = $1 AND host_id = $2`, gameID, hostID)

	return errors.Wrap(err, "cannot delete game")
}

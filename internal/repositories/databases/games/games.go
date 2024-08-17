package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"rollbringer/internal"
)

var gameViewColumns = map[internal.GameView]string{
	internal.GameViewAll: `games.*`,
}

type dbGame struct {
	ID uuid.UUID `db:"id"`

	HostID uuid.UUID `db:"host_id"`
	Name   string    `db:"name"`
}

func (game *dbGame) internalized() *internal.Game {
	if game != nil {
		return &internal.Game{
			ID:     game.ID,
			HostID: game.HostID,
			Name:   game.Name,
		}
	}
	return nil
}

func (db *GamesDatabase) GameInsert(ctx context.Context, game *internal.Game) error {
	err := db.TX.QueryRowxContext(ctx,
		`INSERT INTO games (id, host_id, name)
			VALUES ($1, $2, $3)
		RETURNING id`,
		uuid.New(), game.HostID, game.Name,
	).Scan(&game.ID)

	return errors.Wrap(err, "cannot insert user")
}

func (db *GamesDatabase) GamesCount(ctx context.Context, hostID uuid.UUID) (count int, err error) {
	err = sqlx.GetContext(ctx, db.TX, &count,
		`SELECT COUNT(*) FROM games WHERE host_id = $1`, hostID)

	return count, errors.Wrap(err, "cannot count games")
}

func (db *GamesDatabase) GamesGetForHost(ctx context.Context, hostID uuid.UUID, view internal.GameView) ([]*internal.Game, error) {
	query := fmt.Sprintf(
		`SELECT %s FROM games WHERE games.host_id = $1`,
		gameViewColumns[view],
	)

	var games []*dbGame
	if err := sqlx.SelectContext(ctx, db.TX, &games, query, hostID); err != nil {
		return nil, errors.Wrap(err, "cannot select games")
	}

	// Internalize each game.
	ret := make([]*internal.Game, len(games))
	for i, m := range games {
		ret[i] = m.internalized()
	}

	return ret, nil
}

func (db *GamesDatabase) GameGet(ctx context.Context, gameID uuid.UUID, view internal.GameView) (*internal.Game, error) {
	query := fmt.Sprintf(
		`SELECT %s FROM games WHERE games.id = $1`,
		gameViewColumns[view],
	)

	var game dbGame
	if err := sqlx.GetContext(ctx, db.TX, &game, query, gameID); err != nil {
		if err == sql.ErrNoRows {
			return nil, &internal.ProblemDetail{
				Type:   internal.PDTypeGameNotFound,
				Detail: "Can't find a game with the given game-ID.",
				Extra: map[string]any{
					"game_id": gameID,
				},
			}
		}

		return nil, errors.Wrap(err, "cannot select game")
	}

	return game.internalized(), nil
}

func (db *GamesDatabase) GameDelete(ctx context.Context, gameID, hostID uuid.UUID) error {
	_, err := db.TX.ExecContext(ctx,
		`DELETE FROM games WHERE id = $1 AND host_id = $2`, gameID, hostID)

	return errors.Wrap(err, "cannot delete game")
}

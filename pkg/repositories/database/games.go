package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
)

var gameViewColumns = map[domain.GameView]string{
	domain.GameViewAll:          `games.*`,
	domain.GameViewAll_HostInfo: `games.*, users.id AS "host.id", users.username AS "host.username"`,
}

type gameModel struct {
	ID uuid.UUID `db:"id"`

	HostID uuid.UUID  `db:"host_id"`
	Host   *userModel `db:"host"`

	Name string `db:"name"`
}

func (game *gameModel) domain() *domain.Game {
	if game != nil {
		return &domain.Game{
			ID:     game.ID,
			HostID: game.HostID,
			Host:   game.Host.domain(),
			Name:   game.Name,
		}
	}
	return nil
}

func (db *Database) InsertGame(ctx context.Context, game *domain.Game) error {

	model := gameModel{
		ID:     uuid.New(),
		HostID: game.HostID,
		Name:   game.Name,
	}

	// Insert the game.
	err := sqlx.GetContext(ctx, db.tx, &model,
		`INSERT INTO games (id, host_id, name)
			VALUES ($1, $2, $3)
		RETURNING id`,
		model.ID, model.HostID, model.Name,
	)

	if err != nil {
		return errors.Wrap(err, "cannot insert user")
	}

	*game = *model.domain()
	return nil
}

func (db *Database) GamesCount(ctx context.Context, hostID uuid.UUID) (int, error) {

	// Count games with the host-ID.
	var count int
	err := sqlx.GetContext(ctx, db.tx, &count,
		`SELECT COUNT(*) FROM games WHERE host_id = $1`,
		hostID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, &domain.NormalError{
				Type:   domain.NETypeGameNotFound,
				Detail: "Cannot find a game hosted by a user with the user-ID.",
			}
		}

		return 0, errors.Wrap(err, "cannot count games")
	}

	return count, nil
}

func (db *Database) GetGamesForHost(ctx context.Context, hostID uuid.UUID, view domain.GameView) ([]*domain.Game, error) {

	// Build a query to select games with the host-ID.
	query := fmt.Sprintf(
		`SELECT %s FROM games WHERE games.host_id = $1`,
		gameViewColumns[view],
	)

	// Execute the query.
	var models []*gameModel
	if err := sqlx.SelectContext(ctx, db.tx, &models, query, hostID); err != nil {
		if err == sql.ErrNoRows {
			return nil, &domain.NormalError{
				Type:   domain.NETypeGameNotFound,
				Detail: "Cannot find a game hosted by a user with the user-ID",
			}
		}

		return nil, errors.Wrap(err, "cannot select games")
	}

	// Convert each model to a domain.Game.
	ret := make([]*domain.Game, len(models))
	for i, m := range models {
		ret[i] = m.domain()
	}

	return ret, nil
}

func (db *Database) GetJoinedGamesForUser(ctx context.Context, userID uuid.UUID, view domain.GameView) ([]*domain.Game, error) {

	var joins string
	if view == domain.GameViewAll_HostInfo {
		joins = `LEFT JOIN users ON users.id = games.host_id`
	}

	// Build a query to select games with with the joined user.
	query := fmt.Sprintf(
		`SELECT %s FROM games %s
		WHERE EXISTS (
			SELECT * FROM game_joined_users WHERE game_joined_users.user_id = $1 AND game_joined_users.game_id = games.id
		)`,
		gameViewColumns[view], joins,
	)

	// Execute the query.
	var models []*gameModel
	if err := sqlx.SelectContext(ctx, db.tx, &models, query, userID); err != nil {
		return nil, errors.Wrap(err, "cannot select games")
	}

	// Convert each model to a domain.Game.
	ret := make([]*domain.Game, len(models))
	for i, m := range models {
		ret[i] = m.domain()
	}

	return ret, nil
}

func (db *Database) GetGame(ctx context.Context, gameID uuid.UUID, view domain.GameView) (*domain.Game, error) {

	var joins string

	switch view {
	case domain.GameViewAll_HostInfo:
		joins = `LEFT JOIN users ON users.id = games.host_id`
	}

	// Build a query to select a game with the game-ID.
	query := fmt.Sprintf(
		`SELECT %s FROM games %s WHERE games.id = $1`,
		gameViewColumns[view], joins,
	)

	// Execute the query.
	var model gameModel
	if err := sqlx.GetContext(ctx, db.tx, &model, query, gameID); err != nil {
		if err == sql.ErrNoRows {
			return nil, &domain.NormalError{
				Type:   domain.NETypeGameNotFound,
				Detail: fmt.Sprintf("Cannot find a game with the game-ID"),
			}
		}

		return nil, errors.Wrap(err, "cannot select game")
	}

	return model.domain(), nil
}

func (db *Database) DeleteGame(ctx context.Context, gameID, hostID uuid.UUID) error {

	// Delete the game with the game and host IDs.
	_, err := db.tx.ExecContext(ctx,
		`DELETE FROM games WHERE id = $1 AND host_id = $2`,
		gameID, hostID,
	)

	return errors.Wrap(err, "cannot delete game")
}

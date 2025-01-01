package games

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"rollbringer/internal"
	"rollbringer/internal/repositories/database"
)

func rollColumns(view internal.RollView) string {
	switch view {
	case internal.RollViewListItem:
		return `rolls.*,` +
			`users.id AS "owner.id",` +
			`users.username AS "owner.username",` +
			`users.google_picture AS "owner.google_picture"`

	default:
		return `rolls.*`
	}
}

func rollJoins(view internal.RollView) string {
	switch view {
	case internal.RollViewListItem:
		return `LEFT JOIN users.users ON users.id = rolls.owner_id`
	default:
		return ``
	}
}

func (db *gamesSchema) RollInsert(ctx context.Context, roll *internal.Roll) error {
	query := ` 
		WITH inserted_roll AS (
			INSERT INTO games.rolls (id, owner_id, game_id, dice_types, dice_results, modifiers)
				VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING *
		)
		SELECT inserted_roll.*, users.id AS "owner.id", users.username AS "owner.username", users.google_picture AS "owner.google_picture"
			FROM inserted_roll LEFT JOIN users.users ON users.id = inserted_roll.owner_id
	`

	var dbRoll database.Roll
	err := sqlx.GetContext(ctx, db.TX, &dbRoll, query,
		uuid.New(), roll.OwnerID, roll.GameID, pq.Array(roll.DiceTypes), pq.Array(roll.DiceResults), roll.Modifiers)

	if err != nil {
		return errors.Wrap(err, "cannot insert roll")
	}

	*roll = *dbRoll.Internalized()
	return nil
}

func (db *gamesSchema) RollsGetByGame(ctx context.Context, gameID uuid.UUID, view internal.RollView) ([]*internal.Roll, error) {

	var rolls []*database.Roll
	query := fmt.Sprintf(`SELECT %s FROM games.rolls %s WHERE rolls.game_id = $1`, rollColumns(view), rollJoins(view))
	if err := sqlx.SelectContext(ctx, db.TX, &rolls, query, gameID); err != nil {
		return nil, errors.Wrap(err, "cannot select rolls")
	}

	// Internalize each roll.
	ret := make([]*internal.Roll, len(rolls))
	for i, m := range rolls {
		ret[i] = m.Internalized()
	}

	return ret, nil
}

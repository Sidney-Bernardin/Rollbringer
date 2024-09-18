package games

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"rollbringer/internal"
	"rollbringer/internal/repositories/database"
)

func (db *gamesSchema) RollInsert(ctx context.Context, roll *internal.Roll) error {
	err := db.TX.QueryRowxContext(ctx,
		`INSERT INTO rolls (id, owner_id, game_id, dice_names, dice_results, modifiers)
			VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		uuid.New(), roll.OwnerID, roll.GameID, pq.Array(roll.DiceTypes), pq.Array(roll.DiceResults), roll.Modifiers,
	).Scan(&roll.ID)

	return errors.Wrap(err, "cannot insert roll")
}

func (db *gamesSchema) RollsGetForGame(ctx context.Context, gameID uuid.UUID) ([]*internal.Roll, error) {

	var rolls []*database.Roll
	err := sqlx.SelectContext(ctx, db.TX, &rolls,
		`SELECT rolls.* FROM rolls WHERE rolls.game_id = $1`, gameID)

	if err != nil {
		return nil, errors.Wrap(err, "cannot select rolls")
	}

	// Internalize each roll.
	ret := make([]*internal.Roll, len(rolls))
	for i, m := range rolls {
		ret[i] = m.Internalized()
	}

	return ret, nil
}

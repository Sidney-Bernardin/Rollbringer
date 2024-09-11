package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"rollbringer/internal"
	"rollbringer/internal/repositories/databases"
)

func (db *gamesSchema) RollInsert(ctx context.Context, roll *internal.Roll) error {
	err := db.TX.QueryRowxContext(ctx,
		`INSERT INTO rolls (id, owner_id, game_id, dice_names, dice_results)
			VALUES ($1, $2, $3, $4, $5)
		RETURNING id`,
		uuid.New(), roll.OwnerID, roll.GameID, pq.Array(roll.DiceNames), pq.Array(roll.DiceResults),
	).Scan(&roll.ID)

	return errors.Wrap(err, "cannot insert Roll")
}

func (db *gamesSchema) RollsGetForGame(ctx context.Context, gameID uuid.UUID) ([]*internal.Roll, error) {

	var rolls []*databases.Roll
	err := sqlx.SelectContext(ctx, db.TX, &rolls,
		`SELECT rolls.* FROM rolls WHERE rolls.game_id = $1`, gameID)

	if err != nil {
		return nil, errors.Wrap(err, "cannot select Rolls")
	}

	// Internalize each roll.
	ret := make([]*internal.Roll, len(rolls))
	for i, m := range rolls {
		ret[i] = m.Internalized()
	}

	return ret, nil
}

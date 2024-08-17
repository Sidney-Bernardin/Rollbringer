package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"rollbringer/internal"
)

type dbRoll struct {
	ID uuid.UUID `db:"id"`

	OwnerID uuid.UUID `db:"owner_id"`

	GameID uuid.UUID `db:"game_id"`
	Game   *dbGame   `db:"game"`

	DiceNames   pq.Int32Array `db:"dice_names"`
	DiceResults pq.Int32Array `db:"dice_results"`
}

func (roll *dbRoll) internalized() *internal.Roll {
	if roll != nil {
		return &internal.Roll{
			ID:          roll.ID,
			OwnerID:     roll.OwnerID,
			GameID:      roll.GameID,
			Game:        roll.Game.internalized(),
			DiceNames:   roll.DiceNames,
			DiceResults: roll.DiceResults,
		}
	}
	return nil
}

func (db *GamesDatabase) RollInsert(ctx context.Context, roll *internal.Roll) error {
	err := db.TX.QueryRowxContext(ctx,
		`INSERT INTO rolls (id, owner_id, game_id, dice_names, dice_results)
			VALUES ($1, $2, $3, $4, $5)
		RETURNING id`,
		uuid.New(), roll.OwnerID, roll.GameID, pq.Array(roll.DiceNames), pq.Array(roll.DiceResults),
	).Scan(&roll.ID)

	return errors.Wrap(err, "cannot insert Roll")
}

func (db *GamesDatabase) RollsGetForGame(ctx context.Context, gameID uuid.UUID) ([]*internal.Roll, error) {

	var rolls []*dbRoll
	err := sqlx.SelectContext(ctx, db.TX, &rolls,
		`SELECT rolls.* FROM rolls WHERE rolls.game_id = $1`, gameID)

	if err != nil {
		return nil, errors.Wrap(err, "cannot select Rolls")
	}

	// Internalize each roll.
	ret := make([]*internal.Roll, len(rolls))
	for i, m := range rolls {
		ret[i] = m.internalized()
	}

	return ret, nil
}

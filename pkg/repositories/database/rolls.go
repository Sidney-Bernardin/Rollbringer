package database

import (
	"context"
	"rollbringer/pkg/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

var rollViewColumns = map[domain.RollView]string{
	domain.RollViewAll: `rolls.*`,
}

type rollModel struct {
	ID uuid.UUID `db:"id"`

	OwnerID uuid.UUID  `db:"owner_id"`
	Owner   *userModel `db:"owner"`

	GameID uuid.UUID  `db:"game_id"`
	Game   *gameModel `db:"game"`

	DiceNames   pq.Int32Array `db:"dice_names"`
	DiceResults pq.Int32Array `db:"dice_results"`
}

func (roll *rollModel) domain() *domain.Roll {
	if roll != nil {
		return &domain.Roll{
			ID:          roll.ID,
			OwnerID:     roll.OwnerID,
			Owner:       roll.Owner.domain(),
			DiceNames:   roll.DiceNames,
			DiceResults: roll.DiceResults,
		}
	}
	return nil
}

func (db *Database) InsertRoll(ctx context.Context, roll *domain.Roll) error {

	// Insert the Roll.
	err := db.tx.QueryRowxContext(ctx,
		`INSERT INTO rolls (id, owner_id, game_id, dice_names, dice_results)
			VALUES ($1, $2, $3, $4, $5)
		RETURNING id`,
		uuid.New(), roll.OwnerID, roll.GameID, pq.Array(roll.DiceNames), pq.Array(roll.DiceResults),
	).Scan(&roll.ID)

	return errors.Wrap(err, "cannot insert Roll")
}

func (db *Database) GetRollsForGame(ctx context.Context, gameID uuid.UUID) ([]*domain.Roll, error) {

	// Build a query to select Rolls with the owner-ID.
	query := `SELECT rolls.* FROM rolls WHERE rolls.game_id = $1`

	// Execute the query.
	var models []*rollModel
	if err := sqlx.SelectContext(ctx, db.tx, &models, query, gameID); err != nil {
		return nil, errors.Wrap(err, "cannot select Rolls")
	}

	// Convert each model to a domain.Roll.
	ret := make([]*domain.Roll, len(models))
	for i, m := range models {
		ret[i] = m.domain()
	}

	return ret, nil
}

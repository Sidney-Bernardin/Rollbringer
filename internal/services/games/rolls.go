package games

import (
	"context"
	"rollbringer/internal"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (svc *service) CreateRoll(ctx context.Context, session *internal.Session, gameID uuid.UUID, diceTypes []int, modifiers string) error {
	roll := &internal.Roll{
		OwnerID:     session.UserID,
		GameID:      gameID,
		DiceTypes:   []int32{},
		DiceResults: []int32{},
		Modifiers:   modifiers,
	}

	for _, dType := range diceTypes {
		dType32 := int32(dType)
		roll.DiceTypes = append(roll.DiceTypes, dType32)
		roll.DiceResults = append(roll.DiceResults, svc.random.Int31n(dType32)+1)
	}

	err := svc.schema.RollInsert(ctx, roll)
	return errors.Wrap(err, "cannot insert roll")
}
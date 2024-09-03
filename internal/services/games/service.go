package games

import (
	"context"
	"log/slog"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"rollbringer/internal"
	"rollbringer/internal/config"
	"rollbringer/internal/services"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Service interface {
	services.Servicer
}

type service struct {
	*services.Service

	ps internal.PubSub
	db internal.GamesDatabase

	random *rand.Rand
}

func NewService(
	cfg *config.Config,
	logger *slog.Logger,
	ps internal.PubSub,
	db internal.GamesDatabase,
) Service {
	return &service{
		Service: &services.Service{
			Config: cfg,
			Logger: logger,
		},
		ps:     ps,
		db:     db,
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (svc *service) Shutdown() error {
	svc.ps.Close()
	err := svc.db.Close()
	return errors.Wrap(err, "cannot close database")
}

func (svc *service) getGame(ctx context.Context, gameID uuid.UUID, view internal.GameView) (*internal.Game, error) {
	game, err := svc.db.GameGet(ctx, gameID, view)
	return game, errors.Wrap(err, "cannot get game")
}

func (svc *service) roll(ctx context.Context, diceNamesStr string) (*internal.Roll, error) {

	roll := &internal.Roll{
		DiceNames:   []int32{},
		DiceResults: []int32{},
	}

	for _, dieNameStr := range strings.Split(diceNamesStr, "d")[1:] {
		dName, err := strconv.ParseInt(dieNameStr, 10, 32)
		if err != nil {
			return nil, &internal.ProblemDetail{
				Instance: ctx.Value(internal.CtxKeyInstance).(string),
				Type:     internal.PDTypeInvalidDie,
				Detail:   "Die names must resemble 32-bit integers.",
				Extra: map[string]any{
					"die_name": dieNameStr,
				},
			}
		}

		roll.DiceNames = append(roll.DiceNames, int32(dName))
	}

	for _, dieName := range roll.DiceNames {
		roll.DiceResults = append(roll.DiceResults, svc.random.Int31n(dieName)+1)
	}

	return roll, nil
}

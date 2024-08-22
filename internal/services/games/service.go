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
	database "rollbringer/internal/repositories/databases/games"
	"rollbringer/internal/repositories/pubsub"
	"rollbringer/internal/services"

	"github.com/google/uuid"
)

type Service interface {
	services.Servicer

	Authenticate(ctx context.Context, sessionID uuid.UUID, csrfToken string) (*internal.Session, error)
}

type service struct {
	*services.Service

	ps *pubsub.PubSub
	db *database.GamesDatabase

	random *rand.Rand
}

func NewService(ctx context.Context, cfg *config.Config, logger *slog.Logger, ps *pubsub.PubSub, db *database.GamesDatabase) Service {
	svc := &service{
		Service: &services.Service{
			Config: cfg,
			Logger: logger.With("component", "games_service"),
		},
		ps:     ps,
		db:     db,
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	go svc.doSubscriptions(ctx)
	return svc
}

func (svc *service) Authenticate(ctx context.Context, sessionID uuid.UUID, csrfToken string) (*internal.Session, error) {
	return nil, nil
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

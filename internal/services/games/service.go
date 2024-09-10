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
	services.BaseServicer

	CreateGame(ctx context.Context, session *internal.Session, game *internal.Game) error
	DeleteGame(ctx context.Context, session *internal.Session, gameID uuid.UUID) error
}

type service struct {
	*services.BaseService

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
		BaseService: &services.BaseService{
			Config: cfg,
			Logger: logger,
			PS:     ps,
		},
		db:     db,
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (svc *service) Shutdown() error {
	svc.PS.Close()
	err := svc.db.Close()
	return errors.Wrap(err, "cannot close database")
}

func (svc *service) CreateGame(ctx context.Context, session *internal.Session, game *internal.Game) error {
	return nil
}

func (svc *service) getGame(ctx context.Context, gameID uuid.UUID, view string) (*internal.Game, error) {
	game, err := svc.db.GameGet(ctx, gameID, view)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get game")
	}

	if req.HostView != internal.UserViewNone {
		err := svc.PS.Request(ctx, "users", game.Host, &internal.EventWrapper[any]{
			Event: internal.EventGetUserRequest,
			Payload: internal.GetUserRequest{
				UserID:   game.HostID,
				UserView: req.HostView,
			},
		})

		if err != nil {
			return nil, errors.Wrap(err, "cannot get user")
		}
	}

	return game, nil
}

func (svc *service) DeleteGame(ctx context.Context, session *internal.Session, gameID uuid.UUID) error {
	return nil
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

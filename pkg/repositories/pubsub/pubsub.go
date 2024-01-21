package pubsub

import (
	"context"
	"encoding/json"
	"rollbringer/pkg/models"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type PubSub struct {
	client *redis.Client
	logger *zerolog.Logger
}

// New returns a new PubSub that connects to a Redis server.
func New(loggerCtx context.Context, addr, passw string) (*PubSub, error) {

	// Connect to the Redis server.
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passw,
	})

	// Ping the Redis server.
	if err := c.Ping(context.Background()).Err(); err != nil {
		return nil, errors.Wrap(err, "cannot ping redis server")
	}

	return &PubSub{
		client: c,
		logger: log.Ctx(loggerCtx),
	}, nil
}

func (ps *PubSub) SubToGame(ctx context.Context, gameID uuid.UUID, msgChan chan *models.GameEvent) {

	sub := ps.client.Subscribe(ctx, gameID.String())
	defer sub.Close()

	for {

		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			ps.logger.Error().Stack().Err(err).Msg("Cannot receive message")
			return
		}

		var payload models.GameEvent
		if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
			ps.logger.Error().Stack().Err(err).Msg("Cannot decoded game event")
			return
		}

		msgChan <- &payload
	}
}

func (ps *PubSub) PubToGame(ctx context.Context, gameID uuid.UUID, msg *models.GameEvent) error {

	bytes, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "cannot encode game event")
	}

	err = ps.client.Publish(ctx, gameID.String(), bytes).Err()
	return errors.Wrap(err, "cannot publish game event")
} 

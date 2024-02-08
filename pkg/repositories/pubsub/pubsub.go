package pubsub

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"rollbringer/pkg/domain"
)

type PubSub struct {
	client *redis.Client
	logger *zerolog.Logger
}

// New returns a new PubSub that connects to a Redis server.
func New(rootLogger *zerolog.Logger, addr, passw string) (*PubSub, error) {

	// Connect to the Redis server.
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passw,
	})

	// Ping the Redis server.
	if err := c.Ping(context.Background()).Err(); err != nil {
		return nil, errors.Wrap(err, "cannot ping redis server")
	}

	pubsubLogger := rootLogger.With().Str("component", "pubsub").Logger()

	return &PubSub{
		client: c,
		logger: &pubsubLogger,
	}, nil
}

func (ps *PubSub) SubToGame(ctx context.Context, gameID string, subChan chan *domain.GameEvent) {

	sub := ps.client.Subscribe(ctx, gameID)
	defer sub.Close()

	for {

		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			ps.logger.Error().Stack().Err(err).Msg("Cannot receive message")
			return
		}

		var event *domain.GameEvent
		if err = json.Unmarshal([]byte(msg.Payload), &subChan); err != nil {
			ps.logger.Error().Stack().Err(err).Msg("Cannot decode game event")
			return
		}

		subChan <- event
	}
}

func (ps *PubSub) PubToGame(ctx context.Context, topic string, event *domain.GameEvent) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "cannot encode game event")
	}

	err = ps.client.Publish(ctx, topic, eventBytes).Err()
	return errors.Wrap(err, "cannot publish game event")
}

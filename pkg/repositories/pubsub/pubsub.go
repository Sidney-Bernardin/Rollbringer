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
func New(logger *zerolog.Logger, addr, passw string) (*PubSub, error) {

	// Connect to the Redis server.
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passw,
	})

	// Ping the Redis server.
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, errors.Wrap(err, "cannot ping redis server")
	}

	return &PubSub{client, logger}, nil
}

func (ps *PubSub) Sub(ctx context.Context, topic string, subChan chan domain.Event) {

	sub := ps.client.Subscribe(ctx, topic)
	defer sub.Close()

	for {

		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			ps.logger.Error().Stack().Err(err).Msg("Cannot receive message")
			return
		}

		var event domain.Event
		if err = json.Unmarshal([]byte(msg.Payload), &event); err != nil {
			ps.logger.Error().Stack().Err(err).Msg("Cannot decode event")
			return
		}

		subChan <- event
	}
}

func (ps *PubSub) Pub(ctx context.Context, topic string, data domain.Event) error {

	bytes, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "cannot encode event")
	}

	err = ps.client.Publish(ctx, topic, bytes).Err()
	return errors.Wrap(err, "cannot publish event")
}

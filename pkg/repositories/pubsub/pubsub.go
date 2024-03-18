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

// SubToGame receives events from the game's topic and sends them to subChan.
func (ps *PubSub) SubToGame(ctx context.Context, gameID string, subChan chan domain.Event) {

	// Subscribe to the game's Redis channel.
	sub := ps.client.Subscribe(ctx, gameID)
	defer sub.Close()

	for {

		// Receive the next message.
		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			ps.logger.Error().Stack().Err(err).Msg("Cannot receive message")
			return
		}

		// Decode the event.
		var event domain.Event
		if err = json.Unmarshal([]byte(msg.Payload), &event); err != nil {
			ps.logger.Error().Stack().Err(err).Msg("Cannot decode event")
			return
		}

		subChan <- event
	}
}

// PubToGame sends the event to the game's topic.
func (ps *PubSub) PubToGame(ctx context.Context, gameID string, event domain.Event) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "cannot encode event")
	}

	err = ps.client.Publish(ctx, gameID, eventBytes).Err()
	return errors.Wrap(err, "cannot publish event")
}

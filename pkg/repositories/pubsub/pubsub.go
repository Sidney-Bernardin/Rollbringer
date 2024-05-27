package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"rollbringer/pkg/domain"
)

type PubSub struct {
	client *redis.Client
}

// New returns a new PubSub that connects to a Redis server.
func New(addr, passw string) (*PubSub, error) {

	// Create connection to the Redis server.
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passw,
	})

	// Ping the Redis server.
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, errors.Wrap(err, "cannot ping redis server")
	}

	return &PubSub{client}, nil
}

func (ps *PubSub) Sub(ctx context.Context, topic string, responseChan chan domain.Event, errChan chan error) {

	sub := ps.client.Subscribe(ctx, topic)
	defer sub.Close()

	for {

		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			errChan <- errors.Wrap(err, "cannot receive message")
			return
		}

		event, err := domain.DecodeJSONEvent(ctx, []byte(msg.Payload))
		if err != nil {
			errChan <- fmt.Errorf("cannot decode event from pubsub server: %v", err)
			return
		}

		responseChan <- event
	}
}

func (ps *PubSub) Pub(ctx context.Context, topic string, data any) error {

	bytes, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "cannot encode event")
	}

	err = ps.client.Publish(ctx, topic, string(bytes)).Err()
	return errors.Wrap(err, "cannot publish event")
}

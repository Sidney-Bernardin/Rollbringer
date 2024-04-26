package pubsub

import (
	"context"
	"encoding/json"
	"reflect"

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

func (ps *PubSub) Sub(ctx context.Context, topic string, subChan chan domain.Event, errChan chan error) {

	sub := ps.client.Subscribe(ctx, topic)
	defer sub.Close()

	for {

		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			errChan <- errors.Wrap(err, "cannot receive message")
			return
		}

		var baseEvent domain.BaseEvent
		if err := json.Unmarshal([]byte(msg.Payload), &baseEvent); err != nil {
			continue
		}

		event, ok := domain.OperationEvents[baseEvent.Operation]
		if !ok {
			continue
		}
		event = reflect.New(reflect.TypeOf(event)).Interface().(domain.Event)

		if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
			errChan <- errors.Wrap(err, "cannot decode event")
			continue
		}

		subChan <- event
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

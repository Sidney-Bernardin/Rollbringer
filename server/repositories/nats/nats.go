package nats

import (
	"context"

	"github.com/Sidney-Bernardin/Rollbringer/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"github.com/pkg/errors"
)

type Nats struct {
	conn *nats.Conn

	sessionsKV jetstream.KeyValue
	usersKV    jetstream.KeyValue
}

func New(ctx context.Context, config *server.Config) (*Nats, error) {

	conn, err := nats.Connect(config.NatsUrl)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to server")
	}

	js, err := jetstream.New(conn)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create jetstream instance")
	}

	sessionsKV, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:  "sessions",
		Storage: jetstream.MemoryStorage,
		TTL:     config.SessionTimeout,
	})

	if err != nil {
		return nil, errors.Wrap(err, "cannot create sessions bucket")
	}

	usersKV, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:  "users",
		Storage: jetstream.MemoryStorage,
	})

	if err != nil {
		return nil, errors.Wrap(err, "cannot create users bucket")
	}

	return &Nats{conn, sessionsKV, usersKV}, nil
}

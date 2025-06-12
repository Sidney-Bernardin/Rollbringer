package pubsub

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/Sidney-Bernardin/Rollbringer/server"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

type MsgType string

type PubSub struct {
	log *slog.Logger

	conn *nats.Conn
}

func New(ctx context.Context, config *server.Config, log *slog.Logger) (*PubSub, error) {

	conn, err := nats.Connect(config.NatsUrl)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to server")
	}

	return &PubSub{log, conn}, nil
}

func (ps *PubSub) pub(subject string, msg any, msgType string) error {

	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "cannot JSON encode message")
	}

	natsMsg := nats.NewMsg(subject)
	natsMsg.Header.Set("type", msgType)
	natsMsg.Data = msgJSON

	err = ps.conn.PublishMsg(natsMsg)
	return errors.Wrap(err, "cannot publish message")
}

func (ps *PubSub) sub(ctx context.Context, subject string, cb func(*nats.Msg) error) error {

	sub, err := ps.conn.SubscribeSync(subject)
	if err != nil {
		return errors.Wrap(err, "cannot subscribe")
	}

	for {
		natsMsg, err := sub.NextMsgWithContext(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}

			return errors.Wrap(err, "cannot get next message")
		}

		go func() {
			if err := cb(natsMsg); err != nil {
				ps.log.Log(ctx, slog.LevelError, "Cannot handle incoming message",
					"err", err.Error(),
					"subject", subject)
			}
		}()
	}
}

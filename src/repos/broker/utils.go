package broker

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pkg/errors"

	"rollbringer/src"
)

func (b *broker) publishEvent(subject string, event any) error {

	// JSON encode the event.
	bEvent, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "cannot JSON encode event")
	}

	// Put the event into a new NATS message.
	msg := nats.NewMsg(subject)
	msg.Header.Add("event_type", reflect.TypeOf(event).Elem().Name())
	msg.Data = bEvent

	// Publish the message.
	err = b.conn.PublishMsg(msg)
	return errors.Wrap(err, "cannot publish")
}

func (b *broker) consume(ctx context.Context, consumer jetstream.Consumer, callback func(msg jetstream.Msg)) error {

	c, err := consumer.Consume(callback,
		jetstream.ConsumeErrHandler(func(consumeCtx jetstream.ConsumeContext, err error) {
			b.log.Log(ctx, src.LevelError, "Cannot consume message", "err", err.Error())
		}))

	if err != nil {
		return errors.Wrap(err, "cannot consume")
	}

	// Block until consumer closes or context is done.
	select {
	case <-c.Closed():
		return nil
	case <-ctx.Done():
		c.Stop()
		return nil
	}
}

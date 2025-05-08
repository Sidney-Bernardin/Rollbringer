package broker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pkg/errors"

	"rollbringer/src"
	"rollbringer/src/domain/services/play"
)

type playBroker struct {
	*PublicBroker

	chatStream jetstream.Stream
}

func NewPlayBroker(ctx context.Context, publicBroker *PublicBroker) (play.Broker, error) {

	js, err := jetstream.New(publicBroker.Conn)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create JetStream instance")
	}

	chatStream, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     "chat",
		Subjects: []string{"rooms.*.chat"},
	})

	if err != nil {
		return nil, errors.Wrap(err, "cannot create chat stream")
	}

	return &playBroker{publicBroker, chatStream}, nil
}

func (b *playBroker) Pub(ctx context.Context, event any) bool {
	if b.PublicBroker.Pub(ctx, event) {
		return true
	}

	var subject string
	switch e := event.(type) {
	case *play.EventChat:
		subject = fmt.Sprintf("rooms.%s.chat", e.RoomID)
	default:
		return false
	}

	if err := b.publishEvent(subject, event); err != nil {
		b.log.Log(ctx, src.LevelError, "Cannot publish play event", "err", err.Error())
	}

	return true
}

func (b *playBroker) SubChat(ctx context.Context, roomID uuid.UUID, callback func(event *play.EventChat)) error {

	consumer, err := b.chatStream.CreateConsumer(ctx, jetstream.ConsumerConfig{
		FilterSubject: fmt.Sprintf("rooms.%s.chat", roomID),
	})

	if err != nil {
		return errors.Wrap(err, "cannot create consumer")
	}

	c, err := consumer.Consume(
		func(msg jetstream.Msg) {
			msg.Ack()

			var event *play.EventChat
			if err := json.Unmarshal(msg.Data(), &event); err != nil {
				b.log.Log(ctx, src.LevelError, "Cannot unmarshal message", "err", err.Error())
				return
			}

			callback(event)
		},
		jetstream.ConsumeErrHandler(func(consumeCtx jetstream.ConsumeContext, err error) {
			b.log.Log(ctx, src.LevelError, "Cannot consume message", "err", err.Error())
		}),
	)

	if err != nil {
		return errors.Wrap(err, "cannot consume")
	}

	select {
	case <-c.Closed():
		return errors.Wrap(err, "cannot consume chat")
	case <-ctx.Done():
		c.Stop()
		return nil
	}
}

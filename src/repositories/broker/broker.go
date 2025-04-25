package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"rollbringer/src"
	"rollbringer/src/services"
)

type broker struct {
	log *slog.Logger

	conn       *nats.Conn
	chatStream jetstream.Stream
}

func New(ctx context.Context, config *src.Config, log *slog.Logger) (services.Broker, error) {
	conn, err := nats.Connect(config.NatsURL)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to NATS server")
	}

	js, err := jetstream.New(conn)
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

	return &broker{log, conn, chatStream}, nil
}

func (b *broker) PubChat(event *services.EventChat) error {
	err := b.publishEvent(fmt.Sprintf("rooms.%s.chat", event.RoomID), event)
	return errors.Wrap(err, "cannot publish to chat")
}

func (b *broker) SubRoom(ctx context.Context, roomID src.UUID, callback func(event any)) error {
	errs, ctx := errgroup.WithContext(ctx)

	// Consume chat stream.
	errs.Go(func() error {

		// Create consumer.
		consumer, err := b.chatStream.CreateConsumer(ctx, jetstream.ConsumerConfig{
			FilterSubject: fmt.Sprintf("rooms.%s.chat", roomID),
		})

		if err != nil {
			return errors.Wrap(err, "cannot create consumer")
		}

		// Start consuming.
		c, err := consumer.Consume(
			func(msg jetstream.Msg) {
				msg.Ack()

				var event *services.EventChat
				if err := json.Unmarshal(msg.Data(), &event); err != nil {
					b.log.Log(ctx, src.LevelError, "Cannot unmarshal message", "err", err.Error())
					return
				}

				callback(event)
			},
			jetstream.ConsumeErrHandler(func(_ jetstream.ConsumeContext, err error) {
				b.log.Log(ctx, src.LevelError, "Cannot consume message", "err", err.Error())
			}),
		)

		if err != nil {
			return errors.Wrap(err, "cannot consume")
		}

		// Wait for consumer or context to finish.
		select {
		case <-c.Closed():
			return errors.Wrap(err, "cannot consume chat")
		case <-ctx.Done():
			c.Stop()
			return nil
		}
	})

	// Subscribe to room.
	// errs.Go(func() error {
	// 	sub, err := b.conn.Subscribe(fmt.Sprintf("rooms.%s.chat", roomID), func(msg *nats.Msg) {
	// 		var event any
	// 		switch msg.Header.Get("event_type") {
	// 		case services.EventBoardCreated:
	// 			event = &services.EventBoardCreated{}
	// 		}
	//
	// 		if err := json.Unmarshal(msg.Data, event); err != nil {
	// 			b.log.Log(ctx, src.LevelError, "Cannot unmarshal message", "err", err.Error())
	// 			return
	// 		}
	//
	// 		callback(event)
	// 	})
	//
	// 	if err != nil {
	// 		return errors.Wrap(err, "cannot subscribe")
	// 	}
	//
	// 	<-ctx.Done()
	// 	return errors.Wrap(sub.Unsubscribe(), "cannot unsubscribe")
	// })

	return errors.Wrap(errs.Wait(), "cannot consume room")
}

func (b *broker) publishEvent(subject string, event any) error {
	bEvent, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "cannot marshal event")
	}

	msg := nats.NewMsg(subject)
	msg.Header.Add("event_type", reflect.TypeOf(event).Name())
	msg.Data = bEvent

	err = b.conn.PublishMsg(msg)
	return errors.Wrap(err, "cannot publish")
}

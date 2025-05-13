package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pkg/errors"

	"rollbringer/src"
	"rollbringer/src/domain"
)

type broker struct {
	log *slog.Logger

	Conn *nats.Conn

	chatStream, canvasStream jetstream.Stream
}

func New(ctx context.Context, config *src.Config, log *slog.Logger) (domain.Broker, error) {

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

	canvasStream, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     "canvas",
		Subjects: []string{"boards.*.canvas"},
	})

	if err != nil {
		return nil, errors.Wrap(err, "cannot create chat stream")
	}

	return &broker{log, conn, chatStream, canvasStream}, nil
}

func (b *broker) Pub(ctx context.Context, event any) {
	subjects := []string{}

	switch e := event.(type) {
	case *domain.EventRoomJoined:
		subjects = append(subjects, fmt.Sprintf("rooms.%s", e.RoomID))
	case *domain.EventNewBoard:
		for _, user := range e.Users {
			subjects = append(subjects, fmt.Sprintf("users.%s", user.UserID))
		}
	case *domain.EventChat:
		subjects = append(subjects, fmt.Sprintf("rooms.%s.chat", e.RoomID))
	case *domain.EventUpdateCanvasNode:
		subjects = append(subjects, fmt.Sprintf("boards.%s.canvas", e.BoardID))
	case *domain.EventSaveCanvas:
		subjects = append(subjects, fmt.Sprintf("jobs.canvas-saves"))
	}

	for _, subj := range subjects {
		if err := b.publishEvent(subj, event); err != nil {
			b.log.Log(ctx, src.LevelError, "Cannot publish event", "err", err.Error())
		}
	}
}

func (b *broker) SubUser(ctx context.Context, userID uuid.UUID, callback func(event any)) error {
	sub, err := b.Conn.Subscribe(fmt.Sprintf("users.%s", userID),
		func(msg *nats.Msg) {
			msg.Ack()

			var event any
			switch msg.Header.Get("event_type") {
			case "EventNewBoard":
				event = &domain.EventNewBoard{}
			default:
				b.log.Log(ctx, src.LevelWarn, "Unknown event",
					"subject", msg.Subject,
					"evnet_type", msg.Header.Get("event_type"))
				return
			}

			if err := json.Unmarshal(msg.Data, event); err != nil {
				b.log.Log(ctx, src.LevelError, "Cannot unmarshal message", "err", err.Error())
				return
			}

			callback(event)
		})

	if err != nil {
		return errors.Wrap(err, "cannot subscribe")
	}

	<-ctx.Done()
	return errors.Wrap(sub.Unsubscribe(), "cannot unsubscribe")
}

func (b *broker) SubRoom(ctx context.Context, roomID uuid.UUID, callback func(event any)) error {
	sub, err := b.Conn.Subscribe(fmt.Sprintf("rooms.%s", roomID),
		func(msg *nats.Msg) {
			msg.Ack()

			var event any
			switch msg.Header.Get("event_type") {
			case "EventRoomJoined":
				event = &domain.EventRoomJoined{}
			default:
				b.log.Log(ctx, src.LevelWarn, "Unknown event",
					"subject", msg.Subject,
					"evnet_type", msg.Header.Get("event_type"))
				return
			}

			if err := json.Unmarshal(msg.Data, event); err != nil {
				b.log.Log(ctx, src.LevelError, "Cannot unmarshal message", "err", err.Error())
				return
			}

			callback(event)
		})

	if err != nil {
		return errors.Wrap(err, "cannot subscribe")
	}

	<-ctx.Done()
	return errors.Wrap(sub.Unsubscribe(), "cannot unsubscribe")
}

func (b *broker) SubChat(ctx context.Context, roomID uuid.UUID, callback func(event *domain.EventChat)) error {

	consumer, err := b.chatStream.CreateConsumer(ctx, jetstream.ConsumerConfig{
		FilterSubject: fmt.Sprintf("rooms.%s.chat", roomID),
	})

	if err != nil {
		return errors.Wrap(err, "cannot create consumer")
	}

	c, err := consumer.Consume(
		func(msg jetstream.Msg) {
			msg.Ack()

			var event *domain.EventChat
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

func (b *broker) SubCanvas(ctx context.Context, boardID uuid.UUID, callback func(event any)) error {

	consumer, err := b.canvasStream.CreateConsumer(ctx, jetstream.ConsumerConfig{
		FilterSubject: fmt.Sprintf("boards.%s.canvas", boardID),
	})

	if err != nil {
		return errors.Wrap(err, "cannot create consumer")
	}

	c, err := consumer.Consume(
		func(msg jetstream.Msg) {
			msg.Ack()

			var event any
			switch msg.Headers().Get("event_type") {
			case "EventUpdateCanvasNode":
				event = &domain.EventUpdateCanvasNode{}
			default:
				b.log.Log(ctx, src.LevelWarn, "Unknown event",
					"subject", msg.Subject,
					"evnet_type", msg.Headers().Get("event_type"))
				return
			}

			if err := json.Unmarshal(msg.Data(), event); err != nil {
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
		return errors.Wrap(err, "cannot consume canvas")
	case <-ctx.Done():
		c.Stop()
		return nil
	}
}

func (b *broker) publishEvent(subject string, event any) error {
	bEvent, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "cannot marshal event")
	}

	msg := nats.NewMsg(subject)
	msg.Header.Add("event_type", reflect.TypeOf(event).Elem().Name())
	msg.Data = bEvent

	err = b.Conn.PublishMsg(msg)
	return errors.Wrap(err, "cannot publish")
}

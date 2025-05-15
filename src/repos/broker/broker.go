package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pkg/errors"

	"rollbringer/src"
	"rollbringer/src/domain"
)

type broker struct {
	log *slog.Logger

	conn *nats.Conn

	chatStream,
	canvasesStream jetstream.Stream
}

func New(ctx context.Context, config *src.Config, log *slog.Logger) (domain.Broker, error) {

	// Connect to NATS server.
	conn, err := nats.Connect(config.NatsURL)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to NATS server")
	}

	// Initialize JetStream.
	js, err := jetstream.New(conn)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create JetStream instance")
	}

	// Create the chat stream.
	chatStream, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     "chat",
		Subjects: []string{"rooms.*.chat"},
	})

	if err != nil {
		return nil, errors.Wrap(err, "cannot create chat stream")
	}

	// Create the canvases stream.
	canvasesStream, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name: "canvases",
		Subjects: []string{
			"boards.*.canvas",
			"boards.*.canvas.*",
		},
		NoAck:             true,
		MaxMsgsPerSubject: 5,
	})

	if err != nil {
		return nil, errors.Wrap(err, "cannot create chat stream")
	}

	return &broker{log, conn, chatStream, canvasesStream}, nil
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

	case *domain.EventChatMessage:
		subjects = append(subjects, fmt.Sprintf("rooms.%s.chat", e.RoomID))

	case *domain.EventUpdateCanvasNode:
		subjects = append(subjects, fmt.Sprintf("boards.%s.canvas.%s", e.BoardID, e.Name))
	}

	// Publish the event to each subject.
	for _, subj := range subjects {
		if err := b.publishEvent(subj, event); err != nil {
			b.log.Log(ctx, src.LevelError, "Cannot publish event", "err", err.Error())
		}
	}
}

func (b *broker) SubUser(ctx context.Context, userID uuid.UUID, callback func(event any)) error {
	sub, err := b.conn.Subscribe(fmt.Sprintf("users.%s", userID), func(msg *nats.Msg) {
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

		// JSON decode the message as the event.
		if err := json.Unmarshal(msg.Data, event); err != nil {
			b.log.Log(ctx, src.LevelError, "Cannot JSON decode message", "err", err.Error())
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
	sub, err := b.conn.Subscribe(fmt.Sprintf("rooms.%s", roomID),
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

			// JSON decode the message as the event.
			if err := json.Unmarshal(msg.Data, event); err != nil {
				b.log.Log(ctx, src.LevelError, "Cannot JSON decode message", "err", err.Error())
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

func (b *broker) SubChat(ctx context.Context, roomID uuid.UUID, callback func(event *domain.EventChatMessage)) error {

	consumer, err := b.chatStream.CreateConsumer(ctx, jetstream.ConsumerConfig{
		FilterSubject: fmt.Sprintf("rooms.%s.chat", roomID),
		MaxAckPending: -1,
	})

	if err != nil {
		return errors.Wrap(err, "cannot create consumer")
	}

	err = b.consume(ctx, consumer, func(msg jetstream.Msg) {
		msg.Ack()

		// JSON decode the message as a chat-message event.
		var event *domain.EventChatMessage
		if err := json.Unmarshal(msg.Data(), &event); err != nil {
			b.log.Log(ctx, src.LevelError, "Cannot JSON decode message", "err", err.Error())
			return
		}

		callback(event)
	})

	return errors.Wrap(err, "cannot consume")
}

func (b *broker) SubCanvas(ctx context.Context, boardID uuid.UUID, callback func(event any)) error {

	consumer, err := b.canvasesStream.CreateConsumer(ctx, jetstream.ConsumerConfig{
		FilterSubjects: []string{
			fmt.Sprintf("boards.%s.canvas", boardID),
			fmt.Sprintf("boards.%s.canvas.*", boardID),
		},
		DeliverPolicy: jetstream.DeliverLastPerSubjectPolicy,
		AckPolicy:     jetstream.AckNonePolicy,
	})

	if err != nil {
		return errors.Wrap(err, "cannot create consumer")
	}

	err = b.consume(ctx, consumer, func(msg jetstream.Msg) {

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

		// JSON decode the message as the event.
		if err := json.Unmarshal(msg.Data(), event); err != nil {
			b.log.Log(ctx, src.LevelError, "Cannot JSON decode message", "err", err.Error())
			return
		}

		callback(event)
	})

	return errors.Wrap(err, "cannot consume")
}

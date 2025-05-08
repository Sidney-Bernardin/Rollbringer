package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"

	"rollbringer/src"
	"rollbringer/src/domain"
)

type PublicBroker struct {
	log *slog.Logger

	Conn *nats.Conn
}

func New(ctx context.Context, config *src.Config, log *slog.Logger) (domain.PublicBroker, error) {

	conn, err := nats.Connect(config.NatsURL)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to NATS server")
	}

	return &PublicBroker{log, conn}, nil
}

func (b *PublicBroker) Pub(ctx context.Context, event any) bool {
	subjects := []string{}

	switch e := event.(type) {
	case *domain.EventRoomJoined:
		subjects = append(subjects, fmt.Sprintf("rooms.%s", e.RoomID))
	case *domain.EventNewBoard:
		for _, user := range e.Users {
			subjects = append(subjects, fmt.Sprintf("users.%s", user.UserID))
		}
	default:
		return false
	}

	for _, subj := range subjects {
		if err := b.publishEvent(subj, event); err != nil {
			b.log.Log(ctx, src.LevelError, "Cannot publish public event", "err", err.Error())
		}
	}

	return true
}

func (b *PublicBroker) SubUser(ctx context.Context, userID uuid.UUID, callback func(event any)) error {
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

func (b *PublicBroker) SubRoom(ctx context.Context, roomID uuid.UUID, callback func(event any)) error {
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

func (b *PublicBroker) publishEvent(subject string, event any) error {
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

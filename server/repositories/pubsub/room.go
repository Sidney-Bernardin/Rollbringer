package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Sidney-Bernardin/Rollbringer/server"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/cql"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/sql/queries"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

type (
	ChatMessage struct {
		cql.ChatMessage
		Author *queries.SelectUserRow `json:"author"`
	}
)

const (
	MsgTypeChatMessage MsgType = "chat-message"
)

func (ps *PubSub) PubRoom(roomID server.UUID, msgType MsgType, msg any) error {
	err := ps.pub(fmt.Sprintf("rooms.%s", roomID), msg, string(msgType))
	return errors.WithStack(err)
}

func (ps *PubSub) SubRoom(ctx context.Context, roomID server.UUID, cb func(any)) error {
	err := ps.sub(ctx, fmt.Sprintf("rooms.%s", roomID),
		func(natsMsg *nats.Msg) error {
			var msg any

			switch MsgType(natsMsg.Header.Get("type")) {
			case MsgTypeChatMessage:
				msg = &ChatMessage{}
			}

			if err := json.Unmarshal(natsMsg.Data, msg); err != nil {
				return errors.WithStack(err)
			}

			cb(msg)
			return nil
		})
	return errors.WithStack(err)
}

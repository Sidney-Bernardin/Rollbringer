package service

import (
	"context"
	"log/slog"

	"github.com/Sidney-Bernardin/Rollbringer/server"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/cql"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/pubsub"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/sql/queries"
	"github.com/pkg/errors"
)

func (svc *Service) SendChatMessage(ctx context.Context, roomID server.UUID, author *queries.SelectUserRow, content string) error {
	cqlChatMessage := cql.ChatMessage{
		RoomID:        roomID,
		ChatMessageID: server.NewV1UUID(),
		AuthorID:      author.ID,
		Content:       content,
	}

	if err := svc.CQL.InsertChatMessage(ctx, &cqlChatMessage); err != nil {
		return errors.Wrap(err, "cannot insert chat-message")
	}

	err := svc.PubSub.PubRoom(roomID, pubsub.MsgTypeChatMessage, &pubsub.ChatMessage{
		ChatMessage: cqlChatMessage,
		Author:      author})
	if err != nil {
		svc.Log.Log(ctx, slog.LevelError, "Cannot publish chat-message",
			"err", err.Error(),
			"room_id", roomID,
			"chat_message_id", cqlChatMessage.ChatMessageID,
			"author_id", author.ID)
	}

	return nil
}

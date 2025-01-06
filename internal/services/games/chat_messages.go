package games

import (
	"context"
	"fmt"
	"rollbringer/internal"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (svc *service) CreateChatMessage(ctx context.Context, session *internal.Session, gameID uuid.UUID, message string) (*internal.ChatMessage, error) {
	chatMsg := &internal.ChatMessage{
		OwnerID: session.UserID,
		GameID:  gameID,
		Message: message,
	}

	if err := svc.schema.ChatMessageInsert(ctx, chatMsg); err != nil {
		return nil, errors.Wrap(err, "cannot insert chat-message")
	}

	subject := fmt.Sprintf("games.%s", gameID)
	err := svc.PubSub.Publish(ctx, subject, &internal.EventWrapper[any]{
		Event:   internal.EventChatMessage,
		Payload: chatMsg,
	})

	return chatMsg, errors.Wrap(err, "cannot publish CHAT_MESSAGE event")
}

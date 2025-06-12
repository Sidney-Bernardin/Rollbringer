package cql

import (
	"context"

	"github.com/Sidney-Bernardin/Rollbringer/server"
	"github.com/pkg/errors"
)

type ChatMessage struct {
	RoomID        server.UUID `json:"room_id"`
	ChatMessageID server.UUID `json:"chat_message_id"`
	AuthorID      server.UUID `json:"author_id"`
	Content       string      `json:"content,omitempty"`
}

const qInsertChatMessages = `
	INSERT INTO rollbringer.chat_messages (room_id, chat_message_id, author_id, content)
	VALUES (?, ?, ?, ?)`

func (cql *CQL) InsertChatMessage(ctx context.Context, chatMessage *ChatMessage) error {

	s, err := cql.cluster.CreateSession()
	if err != nil {
		return errors.Wrap(err, "cannot create session")
	}

	err = s.Query(qInsertChatMessages,
		chatMessage.RoomID,
		chatMessage.ChatMessageID,
		chatMessage.AuthorID,
		chatMessage.Content,
	).WithContext(ctx).Exec()
	return errors.Wrap(err, "cannot insert chat-message")
}

const qSelectChatMessages = `SELECT * FROM rollbringer.chat_messages WHERE room_id = ?`

func (cql *CQL) SelectChatMessages(ctx context.Context, roomID server.UUID) ([]*ChatMessage, error) {

	s, err := cql.cluster.CreateSession()
	if err != nil {
		return nil, errors.Wrap(err, "cannot create session")
	}

	scanner := s.Query(qSelectChatMessages, roomID).
		WithContext(ctx).Iter().Scanner()

	chatMessages := []*ChatMessage{}
	for i := 0; scanner.Next(); i++ {

		var chatMsg ChatMessage
		err := scanner.Scan(
			&chatMsg.RoomID,
			&chatMsg.ChatMessageID,
			&chatMsg.AuthorID,
			&chatMsg.Content)
		if err != nil {
			return nil, errors.Wrap(err, "cannot scan row")
		}

		chatMessages = append(chatMessages, &chatMsg)
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.Wrap(err, "cannot select chat-messages")
	}

	return chatMessages, nil
}

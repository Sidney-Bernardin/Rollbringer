package cql

import "github.com/Sidney-Bernardin/Rollbringer/server"

type ChatMessage struct {
	RoomID        server.UUID `json:"room_id"`
	ChatMessageID server.UUID `json:"chat_message_id"`
	AuthorID      server.UUID `json:"author_id"`
	Content       string      `json:"content,omitempty"`
}

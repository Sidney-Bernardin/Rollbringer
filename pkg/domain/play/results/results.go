package results

import "github.com/google/uuid"

type (
	RoomInfo struct {
		RoomID uuid.UUID `json:"room_id"`
		Name   string    `json:"name"`
	}

	RoomListItem struct {
		RoomID  uuid.UUID `json:"room_id"`
		Name    string    `json:"name"`
		Players int       `json:"players"`
	}
)

type (
	ChatMessageInfo struct {
		ChatMessageID uuid.UUID `json:"chat_message_id"`
		AuthorID      uuid.UUID `json:"author_id"`
		Text          string    `json:"text"`
	}
)

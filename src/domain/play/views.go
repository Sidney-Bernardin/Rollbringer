package play

type (
	RoomInfo struct {
		RoomID  string `json:"room_id"`
		OwnerID string `json:"owner_id"`
		Name    string `json:"name"`
	}

	RoomListItem struct {
		RoomID string `json:"room_id"`
		Name   string `json:"name"`
	}
)

type (
	ChatMessageInfo struct {
		ChatMessageID string `json:"chat_message_id"`
		AuthorID      string `json:"author_id"`
		Text          string `json:"text"`
	}
)

package domain

import "github.com/google/uuid"

type (
	EventRoomJoined struct {
		RoomID   uuid.UUID  `json:"room_id"`
		Newcomer PublicUser `json:"newcomer"`
	}

	EventNewBoard struct {
		BoardID uuid.UUID    `json:"board_id"`
		Name    string       `json:"name"`
		Users   []PublicUser `json:"users"`
	}

	EventChat struct {
		RoomID   string `json:"room_id"`
		AuthorID string `json:"author_id"`
		Message  string `json:"message"`
	}

	EventUpdateCanvasNode struct {
		BoardID uuid.UUID `json:"board_id"`
		Name    string    `json:"name"`
		X       int       `json:"x"`
		Y       int       `json:"y"`
	}

	EventSaveCanvas struct {
		BoardID uuid.UUID `json:"board_id"`
	}
)

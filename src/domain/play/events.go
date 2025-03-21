package play

import "github.com/google/uuid"

type (
	EventCanvasUsed struct {
		BoardID uuid.UUID `json:"room_id"`
	}

	EventMovedCanvasNode struct {
		BoardID uuid.UUID `json:"board_id"`
		X       int       `json:"x"`
		Y       int       `json:"y"`
	}
)

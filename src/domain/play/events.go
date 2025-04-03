package play

import "rollbringer/src/domain"

type (
	EventCanvasUsed struct {
		BoardID domain.UUID `json:"room_id"`
	}

	EventMovedCanvasNode struct {
		BoardID domain.UUID `json:"board_id"`
		X       int         `json:"x"`
		Y       int         `json:"y"`
	}
)

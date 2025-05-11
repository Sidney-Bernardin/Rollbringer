package views

import "github.com/google/uuid"

type WebSocketResponse struct {
	Operation string `json:"operation"`
	Payload   any    `json:"payload"`
}

type (
	ReqChat struct {
		Message string `json:"message"`
	}

	ReqGetBoard struct {
		BoardID uuid.UUID `json:"board_id"`
	}
)

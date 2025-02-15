package domain

import "github.com/google/uuid"

type Operation string

const (
	OperationError             Operation = "ERROR"
	OperationSession           Operation = "SESSION"
	OperationGetSessionRequest Operation = "AUTH_REQUEST"
)

type Event struct {
	Operation Operation `json:"operation"`
	Payload   any       `json:"payload"`
}

type GetSessionRequest struct {
	SessionID uuid.UUID `json:"session_id"`
}

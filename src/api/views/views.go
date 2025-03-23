package views

type WSMsgType string

const (
	WSMsgTypeError WSMsgType = "ERROR"
)

type WebSocketMessage struct {
	Type    WSMsgType `json:"type"`
	Payload any       `json:"payload,omitempty"`
}

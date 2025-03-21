package views

type ProblemDetail struct {
	Instance string `json:"instance"`
	Type     string `json:"type"`
	Detail   string `json:"detail,omitempty"`
}

type WSMsgType string

const (
	WSMsgTypeError WSMsgType = "ERROR"
)

type WebSocketMessage struct {
	Type    WSMsgType `json:"type"`
	Payload any       `json:"payload,omitempty"`
}

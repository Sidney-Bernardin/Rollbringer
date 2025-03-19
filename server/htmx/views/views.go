package views

type ProblemDetail struct {
	Instance string `json:"instance"`
	Type     string `json:"type"`
	Detail   string `json:"detail,omitempty"`
}

type WebSocketMessage struct {
	Type    string `json:"type"`
	Payload any    `json:"payload,omitempty"`
}

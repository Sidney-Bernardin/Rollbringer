package views

type WebSocketResponse struct {
	Operation string `json:"operation"`
	Payload   any    `json:"payload"`
}

type (
	ReqChat struct {
		Message string `json:"message"`
	}
)

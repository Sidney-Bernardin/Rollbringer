package domain

type Operation string

const (
	OperationError Operation = "ERROR"
)

type Event struct {
	Operation Operation `json:"operation"`
	Payload   any       `json:"payload"`
}

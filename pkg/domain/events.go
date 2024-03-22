package domain

type Event interface {
	GetOperation() string
}

type BaseEvent struct {
	Headers struct {
		HXTrigger string `json:"hx-trigger"`
	} `json:"HEADERS"`
	Operation string `json:"OPERATION"`
}

func (e BaseEvent) GetOperation() string {
	return e.Operation
}

var OperationTypes = map[string]Event{
	"UPDATE_PDF_FIELD": EventUpdatePDFField{},
	"ERROR":            EventError{},
}

type EventUpdatePDFField struct {
	BaseEvent

	PDFID      string `json:"pdf_id"`
	PageNum    int    `json:"page_num,string"`
	FieldValue string `json:"field_value"`
}

type EventError struct {
	BaseEvent

	Err error `json:"err,omitempty"`
}

func NewEventError(err error) *EventError {
	return &EventError{
		BaseEvent: BaseEvent{
			Operation: "ERROR",
		},
		Err: err,
	}
}

func (e *EventError) Error() string {
	return e.Err.Error()
}

func (e *EventError) Cause() error {
	return e.Err
}

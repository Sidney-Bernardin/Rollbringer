package domain

import (
	"github.com/google/uuid"
)

type Event interface {
	GetHeaders() map[string]any
}

// =====

type BaseEvent struct {
	Headers   map[string]any `json:"HEADERS"`
	Operation string         `json:"OPERATION"`
}

func (e BaseEvent) GetHeaders() map[string]any {
	return e.Headers
}

// =====

var OperationEvents = map[string]Event{
	"SUB_TO_PDF":       EventSubToPDFPage{},
	"UPDATE_PDF_FIELD": EventUpdatePDFField{},
	"ERROR":            EventError{},
}

type EventSubToPDFPage struct {
	BaseEvent

	PDFID uuid.UUID `json:"pdf_id"`
}

type EventUpdatePDFField struct {
	BaseEvent

	PDFID      uuid.UUID `json:"pdf_id"`
	PageNum    int       `json:"page_num,string"`
	FieldName  string    `json:"field_name"`
	FieldValue string    `json:"field_value"`
}

type EventError struct {
	BaseEvent

	Err error `json:"err,omitempty"`
}

func NewEventError(err error) error {
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

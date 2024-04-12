package domain

import (
	"fmt"
	"reflect"

	"github.com/google/uuid"
)

type Event interface {
	GetOperationStruct() (Event, error)
}

type BaseEvent struct {
	Operation string `json:"OPERATION"`
}

var operationToStruct = map[string]Event{
	"SUB_TO_PDF":       EventSubToPDFPage{},
	"UPDATE_PDF_FIELD": EventUpdatePDFField{},
	"ERROR":            EventError{},
}

func (e BaseEvent) GetOperationStruct() (Event, error) {

	event, ok := operationToStruct[e.Operation]
	if !ok {
		return nil, &ProblemDetail{
			Type:   PDTypeInvalidEventOperation,
			Detail: fmt.Sprintf(`"%s" is an invlid event operation`, e.Operation),
		}
	}

	event = reflect.New(reflect.TypeOf(event)).Interface().(Event)
	return event, nil
}

type (
	EventSubToPDFPage struct {
		BaseEvent

		PDFID uuid.UUID `json:"pdf_id"`
	}

	EventUpdatePDFField struct {
		BaseEvent

		PDFID      uuid.UUID `json:"pdf_id"`
		PageNum    int       `json:"page_num,string"`
		FieldName  string    `json:"field_name"`
		FieldValue string    `json:"field_value"`
	}

	EventError struct {
		BaseEvent

		Err error `json:"err,omitempty"`
	}
)

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

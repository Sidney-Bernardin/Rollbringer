package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	OperationNormalError    = "NORMAL_ERROR"
	OperationPDFFields      = "PDF_FIELDS"
	OperationSubToPDF       = "SUB_TO_PDF"
	OperationUpdatePDFField = "UPDATE_PDF_FIELD"
	OperationCreateRoll     = "CREATE_ROLL"
	OperationRoll           = "ROLL"
)

type Event interface {
	GetHeaders() map[string]any
	Validate(ctx context.Context) error
}

type BaseEvent struct {
	Headers   map[string]any `json:"HEADERS"`
	Operation string         `json:"OPERATION"`
}

func (e BaseEvent) GetHeaders() map[string]any {
	return e.Headers
}

func DecodeJSONEvent(ctx context.Context, eventBytes []byte) (Event, error) {

	var baseEvent BaseEvent
	if err := json.Unmarshal(eventBytes, &baseEvent); err != nil {
		return nil, &NormalError{
			Instance: ctx.Value(CtxKeyInstance).(string),
			Type:     NETypeCannotDecodeRequest,
			Detail:   fmt.Sprintf("Your request isn't valid JSON: '%v'", err),
		}
	}

	var event Event
	switch baseEvent.Operation {
	case OperationNormalError:
		event = &EventNormalError{}
	case OperationPDFFields:
		event = &EventPDFFields{}
	case OperationSubToPDF:
		event = &EventSubToPDF{}
	case OperationUpdatePDFField:
		event = &EventUpdatePDFField{}
	case OperationCreateRoll:
		event = &EventCreateRoll{}
	case OperationRoll:
		event = &EventRoll{}
	default:
		return nil, &NormalError{
			Instance: ctx.Value(CtxKeyInstance).(string),
			Type:     NETypeInvalidEventOperation,
			Detail:   fmt.Sprintf("'%s' is an invalid operation.", baseEvent.Operation),
		}
	}

	if err := json.Unmarshal(eventBytes, &event); err != nil {
		return nil, &NormalError{
			Instance: ctx.Value(CtxKeyInstance).(string),
			Type:     NETypeCannotDecodeRequest,
			Detail:   fmt.Sprintf("Your request isn't valid JSON: '%v'", err),
		}
	}

	return event, nil
}

type EventPDFFields struct {
	BaseEvent

	PDFID   uuid.UUID         `json:"pdf_id,omitempty"`
	PageNum int               `json:"page_num,string,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func (event *EventPDFFields) Validate(ctx context.Context) error {
	if event.PageNum < 1 {
		return &NormalError{
			Instance: ctx.Value(CtxKeyInstance).(string),
			Type:     NETypeInvalidPDFPageNumber,
			Detail:   "Page number cannot be less than 1.",
		}
	}

	return nil
}

type EventSubToPDF struct {
	BaseEvent

	PDFID   uuid.UUID `json:"pdf_id,omitempty"`
	PageNum int       `json:"page_num,string,omitempty"`
}

func (event *EventSubToPDF) Validate(ctx context.Context) error {
	if event.PageNum < 1 {
		return &NormalError{
			Instance: ctx.Value(CtxKeyInstance).(string),
			Type:     NETypeInvalidPDFPageNumber,
			Detail:   "Page number cannot be less than 1.",
		}
	}

	return nil
}

type EventUpdatePDFField struct {
	BaseEvent

	PDFID      uuid.UUID `json:"pdf_id,omitempty"`
	PageNum    int       `json:"page_num,string,omitempty"`
	FieldName  string    `json:"field_name,omitempty"`
	FieldValue string    `json:"field_value,omitempty"`
}

func (event *EventUpdatePDFField) Validate(ctx context.Context) error {
	if event.PageNum < 1 {
		return &NormalError{
			Instance: ctx.Value(CtxKeyInstance).(string),
			Type:     NETypeInvalidPDFPageNumber,
			Detail:   "Page number is less than 1.",
		}
	}

	if event.FieldName == "" {
		return &NormalError{
			Instance: ctx.Value(CtxKeyInstance).(string),
			Type:     NETypeInvalidPDFFieldName,
			Detail:   "Field name cannot be empty.",
		}
	}

	if strings.Contains(event.FieldName, " ") {
		return &NormalError{
			Instance: ctx.Value(CtxKeyInstance).(string),
			Type:     NETypeInvalidPDFFieldName,
			Detail:   "Field name cannot have spaces.",
		}
	}

	return nil
}

type EventCreateRoll struct {
	BaseEvent

	Dice     string `json:"dice,omitempty"`
	Modifier string `json:"modifier,omitempty"`
}

func (*EventCreateRoll) Validate(ctx context.Context) error {
	return nil
}

type EventRoll struct {
	BaseEvent
	*Roll
}

func (event *EventRoll) Validate(ctx context.Context) error {
	if err := event.Roll.validate(ctx); err != nil {
		return errors.Wrap(err, "invalid roll")
	}

	return nil
}

type EventNormalError struct {
	BaseEvent
	*NormalError
}

func (*EventNormalError) Validate(ctx context.Context) error {
	return nil
}

package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Operation string

const (
	OperationNormalError    Operation = "NORMAL_ERROR"
	OperationPDFFields      Operation = "PDF_FIELDS"
	OperationSubToPDF       Operation = "SUB_TO_PDF"
	OperationUpdatePDFField Operation = "UPDATE_PDF_FIELD"
)

type Event interface {
	GetHeaders() map[string]any
	Validate(ctx context.Context) error
}

type BaseEvent struct {
	Headers   map[string]any `json:"HEADERS"`
	Operation Operation      `json:"OPERATION"`
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

	PDFID   uuid.UUID         `json:"pdf_id"`
	PageNum int               `json:"page_num,string"`
	Fields  map[string]string `json:"fields"`
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

	PDFID   uuid.UUID `json:"pdf_id"`
	PageNum int       `json:"page_num,string"`
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

	PDFID      uuid.UUID `json:"pdf_id"`
	PageNum    int       `json:"page_num,string"`
	FieldName  string    `json:"field_name"`
	FieldValue string    `json:"field_value"`
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

type EventNormalError struct {
	BaseEvent
	*NormalError
}

func (*EventNormalError) Validate(ctx context.Context) error {
	return nil
}

package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type eventType string

const (
	ETError          eventType = "ERROR"
	ETSubToPDF       eventType = "SUB_TO_PDF"
	ETPdfFields      eventType = "PDF_FIELDS"
	ETUpdatePDFField eventType = "UPDATE_PDF_FIELD"
	ETCreateRoll     eventType = "CREATE_ROLL"
	ETGame           eventType = "GAME"
	ETRoll           eventType = "ROLL"
)

type Event interface {
	GetHeaders() map[string]any
	Validate(ctx context.Context) error
}

type BaseEvent struct {
	Headers map[string]any `json:"HEADERS"`
	Type    eventType      `json:"TYPE"`
}

func (e BaseEvent) GetHeaders() map[string]any {
	return e.Headers
}

func JSONEncodeEvent(ctx context.Context, event Event) ([]byte, error) {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return nil, NewProblemDetail(ctx, &PDOptions{
			Type:   PDTypeInvalidEvent,
			Detail: err.Error(),
		})
	}
	return eventBytes, nil
}

func JSONDecodeEvent(ctx context.Context, eventBytes []byte) (Event, error) {

	var baseEvent BaseEvent
	if err := json.Unmarshal(eventBytes, &baseEvent); err != nil {
		return nil, NewProblemDetail(ctx, &PDOptions{
			Type:   PDTypeInvalidEvent,
			Detail: err.Error(),
		})
	}

	var event Event
	switch baseEvent.Type {
	case ETError:
		event = &EventError{}
	case ETPdfFields:
		event = &EventPDFFields{}
	case ETSubToPDF:
		event = &EventSubToPDF{}
	case ETUpdatePDFField:
		event = &EventUpdatePDFField{}
	case ETCreateRoll:
		event = &EventCreateRoll{}
	case ETRoll:
		event = &EventRoll{}
	default:
		return nil, NewProblemDetail(ctx, &PDOptions{
			Type:   PDTypeInvalidEvent,
			Detail: fmt.Sprintf("'%s' is an invalid event-type.", baseEvent.Type),
			Extra: map[string]any{
				"event_type": baseEvent.Type,
			},
		})
	}

	if err := json.Unmarshal(eventBytes, &event); err != nil {
		return nil, NewProblemDetail(ctx, &PDOptions{
			Type:   PDTypeInvalidEvent,
			Detail: err.Error(),
		})
	}

	return event, nil
}

type EventError struct {
	BaseEvent
	*ProblemDetail
}

func (*EventError) Validate(ctx context.Context) error {
	return nil
}

type EventSubToPDF struct {
	BaseEvent

	PDFID   uuid.UUID `json:"pdf_id,omitempty"`
	PageNum int       `json:"page_num,string,omitempty"`
}

func (event *EventSubToPDF) Validate(ctx context.Context) error {
	if event.PageNum < 1 {
		return &ProblemDetail{
			Instance: ctx.Value(CtxKeyInstance).(string),
			Type:     PDTypeInvalidPDFPageNumber,
			Detail:   "page_num cannot be less than 1.",
		}
	}

	return nil
}

type EventPDFFields struct {
	BaseEvent

	PDFID   uuid.UUID         `json:"pdf_id,omitempty"`
	PageNum int               `json:"page_num,string,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func (event *EventPDFFields) Validate(ctx context.Context) error {
	if event.PageNum < 1 {
		return &ProblemDetail{
			Instance: ctx.Value(CtxKeyInstance).(string),
			Type:     PDTypeInvalidPDFPageNumber,
			Detail:   "page_num cannot be less than 1.",
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
		return &ProblemDetail{
			Instance: ctx.Value(CtxKeyInstance).(string),
			Type:     PDTypeInvalidPDFPageNumber,
			Detail:   "page_num cannot be less than 1.",
		}
	}

	if event.FieldName == "" {
		return &ProblemDetail{
			Instance: ctx.Value(CtxKeyInstance).(string),
			Type:     PDTypeInvalidPDFFieldName,
			Detail:   "field_name cannot be empty.",
		}
	}

	if strings.Contains(event.FieldName, " ") {
		return &ProblemDetail{
			Instance: ctx.Value(CtxKeyInstance).(string),
			Type:     PDTypeInvalidPDFFieldName,
			Detail:   "field_name cannot have spaces.",
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

type EventGame struct {
	BaseEvent
	*Game
}

func (event *EventGame) Validate(ctx context.Context) error {
	return nil
}

type EventRoll struct {
	BaseEvent
	*Roll
}

func (event *EventRoll) Validate(ctx context.Context) error {
	return nil
}

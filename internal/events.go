package internal

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

type Event interface {
	GetHeaders() map[string]any
	Validate(ctx context.Context) error
}

type BaseEvent struct {
	Headers map[string]any `json:"HEADERS"`
	Type    eventType      `json:"TYPE"`
}

func (e BaseEvent) Validate(ctx context.Context) error { return nil }

func (e BaseEvent) GetHeaders() map[string]any {
	return e.Headers
}

type eventType string

const (
	ETError        eventType = "ERROR"
	ETAuthenticate eventType = "AUTHENTICATE"

	ETUser    eventType = "USER"
	ETGetUser eventType = "GET_USER"

	ETSession    eventType = "SESSION"
	ETGetSession eventType = "GET_SESSION"

	ETGame    eventType = "GAME"
	ETGetGame eventType = "GET_GAME"

	ETSubToPDF       eventType = "SUB_TO_PDF"
	ETPdfFields      eventType = "PDF_FIELDS"
	ETUpdatePDFField eventType = "UPDATE_PDF_FIELD"

	ETCreateRoll eventType = "CREATE_ROLL"
	ETRoll       eventType = "ROLL"
)

var eventTypes = map[eventType]Event{
	ETError:        &EventError{},
	ETAuthenticate: &EventAuthenticate{},

	ETUser:    &EventUser{},
	ETGetUser: &EventGetUser{},

	ETSession:    &EventSession{},
	ETGetSession: &EventGetSession{},

	ETGame:    &EventGame{},
	ETGetGame: &EventGetGame{},

	ETSubToPDF:       &EventSubToPDF{},
	ETPdfFields:      &EventPDFFields{},
	ETUpdatePDFField: &EventUpdatePDFField{},

	ETCreateRoll: &EventCreateRoll{},
	ETRoll:       &EventRoll{},
}

type EventError struct {
	BaseEvent
	ProblemDetail
}

type EventAuthenticate struct {
	BaseEvent

	SessionID uuid.UUID `json:"session_id,omitempty"`
	CSRFToken string    `json:"csrf_token,omitempty"`
}

type EventUser struct {
	BaseEvent
	User
}

type EventGetUser struct {
	BaseEvent

	UserID uuid.UUID `json:"user_id,omitempty"`
	View   UserView  `json:"view,omitempty"`
}

type EventSession struct {
	BaseEvent
	Session
}

type EventGetSession struct {
	BaseEvent

	View      SessionView `json:"view,omitempty"`
	SessionID uuid.UUID   `json:"session_id,omitempty"`
}

type EventGame struct {
	BaseEvent
	Game
}

type EventGetGame struct {
	BaseEvent

	GameID uuid.UUID `json:"game_id,omitempty"`
	View   GameView  `json:"view,omitempty"`
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

type EventRoll struct {
	BaseEvent
	Roll
}

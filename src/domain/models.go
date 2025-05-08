package domain

import (
	"fmt"
)

type ExternalErrorType string

const (
	ExternalErrorTypeInternalError ExternalErrorType = "internal-error"
	ExternalErrorTypeUnauthorized  ExternalErrorType = "unauthorized"
	ExternalErrorTypeInvalidUUID   ExternalErrorType = "invalid-uuid"
)

type ExternalError struct {
	Type    ExternalErrorType `json:"type"`
	Msg     string            `json:"description,omitempty"`
	Details map[string]any    `json:"attrs,omitempty"`
}

func (err *ExternalError) Error() string {
	return fmt.Sprintf("%s: %s", err.Type, err.Msg)
}

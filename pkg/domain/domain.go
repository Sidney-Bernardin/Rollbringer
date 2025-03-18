package domain

import (
	"errors"
	"fmt"
	"log/slog"
)

var (
	SlogLevelTrace slog.Level = -8
	LevelDebug     slog.Level = -4
	LevelInfo      slog.Level = 0
	LevelWarn      slog.Level = 4
	LevelError     slog.Level = 8
	SlogLevelFatal slog.Level = 12
)

var (
	ErrEntityConflict = errors.New("entity conflict")
)

type DomainErrorType string

type DomainError struct {
	Type        DomainErrorType
	Description string
	Details     map[string]any
}

func (err *DomainError) Error() string {
	return fmt.Sprintf("%s: %s", err.Type, err.Description)
}

package server

import (
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

type Config struct {
	DevMode        bool `default:"true" split_words:"true"`
	LogPretty      bool `default:"true" split_words:"true"`
	LogPrettyColor bool `default:"true" split_words:"true"`

	HTMXAddr string `required:"true" split_words:"true"`

	PostgresURL string `required:"true" split_words:"true"`
}

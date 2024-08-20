package services

import (
	"log/slog"

	"rollbringer/internal/config"
)

type Servicer interface{}

type Service struct {
	Config *config.Config
	Logger *slog.Logger
}

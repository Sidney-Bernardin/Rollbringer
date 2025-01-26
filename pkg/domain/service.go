package domain

import (
	"context"
	"log/slog"
)

var (
	SlogLevelTrace slog.Level = -8
	SlogLevelFatal slog.Level = 12
)

type IService interface {
	Run(context.Context) error
	Shutdown(context.Context) error
}

type Service struct {
	Config *Config
	Logger *slog.Logger

	PubSub PubSubRepository
}

func (svc *Service) Run(context.Context) error { return nil }

func (svc *Service) Shutdown(ctx context.Context) error {
	svc.PubSub.Close()
	return nil
}

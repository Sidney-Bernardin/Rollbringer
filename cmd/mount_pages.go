//go:build all || pages
// +build all pages

package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/pkg/errors"

	"rollbringer/internal/config"
	handler "rollbringer/internal/handlers/pages"
	"rollbringer/internal/repositories/pubsub"
	service "rollbringer/internal/services/pages"
)

func init() {
	serviceHandlers["/"] = func(ctx context.Context, cfg *config.Config, logger *slog.Logger) (http.Handler, error) {
		ps, err := pubsub.New(cfg, logger)
		if err != nil {
			return nil, errors.Wrap(err, "cannot create pubsub repository")
		}

		svc := service.NewService(cfg, logger, ps)
		return handler.NewHandler(cfg, logger, svc), nil
	}
}

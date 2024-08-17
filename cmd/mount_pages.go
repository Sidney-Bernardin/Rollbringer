//go:build all || pages
// +build all pages

package main

import (
	"log/slog"
	"net/http"

	"github.com/pkg/errors"

	"rollbringer/internal/config"
	handler "rollbringer/internal/handlers/users"
	"rollbringer/internal/repositories/pubsub"
	service "rollbringer/internal/services/users"
)

func init() {
	serviceHandlers["/"] = func(cfg *config.Config, logger *slog.Logger) (http.Handler, error) {
		ps, err := pubsub.New(cfg, logger)
		if err != nil {
			return nil, errors.Wrap(err, "cannot create pubsub repository")
		}

		svc := service.New(cfg, logger, ps)
		return handler.New(cfg, logger, svc), nil
	}
}

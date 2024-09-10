//go:build all || pages
// +build all pages

package main

import (
	"log/slog"
	"net/http"

	"github.com/pkg/errors"

	"rollbringer/internal/config"
	handler "rollbringer/internal/handlers/pages"
	"rollbringer/internal/repositories/pubsub"
	"rollbringer/internal/services"
	service "rollbringer/internal/services/pages"
)

func init() {
	features["pages"] = func(cfg *config.Config, logger *slog.Logger) (http.Handler, services.BaseServicer, error) {

		// Create a PubSub repository.
		ps, err := pubsub.New(cfg, logger)
		if err != nil {
			return nil, nil, errors.Wrap(err, "cannot create pubsub repository")
		}

		// Create a service and handler.
		svc := service.NewService(cfg, logger, ps)
		return handler.NewHandler(cfg, logger, svc), svc, nil
	}
}

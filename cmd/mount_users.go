//go:build all || users
// +build all users

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
	serviceHandlers["/users"] = func(cfg *config.Config, logger *slog.Logger) (http.Handler, error) {
		ps, err := pubsub.New(cfg, logger)
		if err != nil {
			return nil, errors.Wrap(err, "cannot create pubsub repository")
		}

		svc := service.NewService(cfg, logger, ps)
		return handler.NewHandler(cfg, logger, svc), nil
	}
}

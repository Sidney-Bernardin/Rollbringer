//go:build all || games
// +build all games

package main

import (
	"log/slog"
	"net/http"

	"github.com/pkg/errors"

	"rollbringer/internal/config"
	handler "rollbringer/internal/handlers/games"
	database "rollbringer/internal/repositories/databases/games"
	"rollbringer/internal/repositories/pubsub"
	"rollbringer/internal/services"
	service "rollbringer/internal/services/games"
)

func init() {
	features["games"] = func(cfg *config.Config, logger *slog.Logger) (http.Handler, services.Servicer, error) {

		// Create a PubSub repository.
		ps, err := pubsub.New(cfg, logger)
		if err != nil {
			return nil, nil, errors.Wrap(err, "cannot create pubsub repository")
		}

		// Create a GamesDatabase repository.
		db, err := database.New(cfg, logger)
		if err != nil {
			return nil, nil, errors.Wrap(err, "cannot create database repository")
		}

		// Create a service and handler.
		svc := service.NewService(cfg, logger, ps, db)
		return handler.NewHandler(cfg, logger, svc), svc, nil
	}
}

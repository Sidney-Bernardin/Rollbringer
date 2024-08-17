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
	service "rollbringer/internal/services/games"
)

func init() {
	serviceHandlers["/games"] = func(cfg *config.Config, logger *slog.Logger) (http.Handler, error) {
		ps, err := pubsub.New(cfg, logger)
		if err != nil {
			return nil, errors.Wrap(err, "cannot create pubsub repository")
		}

		db, err := database.New(cfg, logger)
		if err != nil {
			return nil, errors.Wrap(err, "cannot create database repository")
		}

		svc := service.New(cfg, logger, ps, db)
		return handler.New(cfg, logger, svc), nil
	}
}

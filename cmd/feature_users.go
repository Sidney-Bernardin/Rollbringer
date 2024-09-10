//go:build all || users
// +build all users

package main

import (
	"log/slog"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"rollbringer/internal/config"
	handler "rollbringer/internal/handlers/users"
	database "rollbringer/internal/repositories/databases/users"
	"rollbringer/internal/repositories/oauth"
	"rollbringer/internal/repositories/pubsub"
	"rollbringer/internal/services"
	service "rollbringer/internal/services/users"
)

func init() {
	features["users"] = func(cfg *config.Config, logger *slog.Logger) (http.Handler, services.BaseServicer, error) {

		// Create a PubSub repository.
		ps, err := pubsub.New(cfg, logger)
		if err != nil {
			return nil, nil, errors.Wrap(err, "cannot create pubsub repository")
		}

		// Create a UsersDatabase repository.
		db, err := database.New(cfg, logger)
		if err != nil {
			return nil, nil, errors.Wrap(err, "cannot create database repository")
		}

		// Create an OAuth repository.
		oa := oauth.New(cfg)

		// Create a service and handler.
		svc := service.NewService(cfg, logger, ps, db, oa)
		return handler.NewHandler(cfg, logger, svc), svc, nil
	}
}

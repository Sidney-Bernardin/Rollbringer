//go:build all || games
// +build all games

package main

import (
	"net/http"

	"github.com/pkg/errors"

	handler "rollbringer/internal/handlers/games"
	schema "rollbringer/internal/repositories/database/games"
	"rollbringer/internal/repositories/pubsub"
	"rollbringer/internal/services"
	service "rollbringer/internal/services/games"
)

func init() {
	features["games"] = func(deps globalDependencies) (http.Handler, services.BaseServicer, error) {

		// Create a PubSub repository.
		pubsubRepo, err := pubsub.New(deps.cfg, deps.logger)
		if err != nil {
			return nil, nil, errors.Wrap(err, "cannot create pubsub repository")
		}

		// Create a GamesSchema repository.
		schemaRepo := schema.New(deps.dbRepo)

		// Create a service and handler.
		svc := service.NewService(deps.cfg, deps.logger, pubsubRepo, schemaRepo)
		return handler.NewHandler(deps.cfg, deps.logger, svc), svc, nil
	}
}

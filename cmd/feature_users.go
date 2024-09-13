//go:build all || users
// +build all users

package main

import (
	"net/http"

	"github.com/pkg/errors"

	handler "rollbringer/internal/handlers/users"
	schema "rollbringer/internal/repositories/database/users"
	"rollbringer/internal/repositories/oauth"
	"rollbringer/internal/repositories/pubsub"
	"rollbringer/internal/services"
	service "rollbringer/internal/services/users"
)

func init() {
	features["users"] = func(deps globalDependencies) (http.Handler, services.BaseServicer, error) {

		// Create a PubSub repository.
		pubsubRepo, err := pubsub.New(deps.cfg, deps.logger)
		if err != nil {
			return nil, nil, errors.Wrap(err, "cannot create pubsub repository")
		}

		// Create a UsersDatabase repository.
		schemaRepo := schema.New(deps.dbRepo)

		// Create an OAuth repository.
		oa := oauth.New(deps.cfg)

		// Create a service and handler.
		svc := service.NewService(deps.cfg, deps.logger, pubsubRepo, schemaRepo, oa)
		return handler.NewHandler(deps.cfg, deps.logger, svc), svc, nil
	}
}

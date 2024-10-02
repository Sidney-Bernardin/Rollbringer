//go:build all || pages
// +build all pages

package main

import (
	"net/http"

	"github.com/pkg/errors"

	handler "rollbringer/internal/handlers/pages"
	"rollbringer/internal/repositories/pubsub"
	"rollbringer/internal/services"
	service "rollbringer/internal/services/pages"
)

func init() {
	features["pages"] = func(deps globalDependencies) (http.Handler, services.BaseServicer, error) {

		// Create a PubSub repository.
		pubsubRepo, err := pubsub.New(deps.cfg, deps.logger)
		if err != nil {
			return nil, nil, errors.Wrap(err, "cannot create pubsub repository")
		}

		// Create a service and handler.
		svc := service.NewService(deps.cfg, deps.logger, pubsubRepo)
		return handler.NewHandler(deps.cfg, deps.logger, svc), svc, nil
	}
}

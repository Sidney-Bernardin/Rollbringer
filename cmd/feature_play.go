//go:build !noplay
// +build !noplay

package main

import (
	"rollbringer/pkg/domain"
	service "rollbringer/pkg/domain/services/play"
	handler "rollbringer/pkg/handlers/play"
	database "rollbringer/pkg/repositories/database/play"
	"rollbringer/pkg/repositories/pubsub"
)

func init() {
	registeredFeatures["play"] = func() error {

		// Create PubSub repository.
		pubSubRepo, err := pubsub.NewPubSubRepository(config, logger.With("dependency", "nats-pubsub-repo"))
		if err != nil {
			return domain.Wrap(err, "cannot create PubSub repository", nil)
		}

		// Create play database repository.
		playDBRepo, err := database.NewPlayDatabaseRepository(config, logger.With("dependency", "db-repo"), migrations)
		if err != nil {
			return domain.Wrap(err, "cannot create accounts database repository", nil)
		}

		svc := service.New(config, logger.With("dependency", "domain"), pubSubRepo, playDBRepo)
		h := handler.New(config, logger.With("dependency", "http-api"), svc)

		features["play"] = feature{h, svc, "/services/play"}
		return nil
	}
}

//go:build !nopages
// +build !nopages

package main

import (
	"rollbringer/pkg/domain"
	"rollbringer/pkg/handlers/pages"
	"rollbringer/pkg/repositories/pubsub"
)

func init() {
	registeredFeatures["pages"] = func() error {

		// Create PubSub repository.
		pubSubRepo, err := pubsub.NewPubSubRepository(config, logger.With("dependency", "nats-pubsub-repo"))
		if err != nil {
			return domain.Wrap(err, "cannot create PubSub repository", nil)
		}

		svc := &domain.Service{
			Config: config,
			Logger: logger.With("dependency", "base-service"),
			PubSub: pubSubRepo,
		}
		h := handler.New(config, logger.With("dependency", "http-api"), svc)

		features["pages"] = feature{h, svc}
		return nil
	}
}

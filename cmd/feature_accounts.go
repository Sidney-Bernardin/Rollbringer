//go:build !noaccounts
// +build !noaccounts

package main

import (
	"rollbringer/pkg/domain"
	service "rollbringer/pkg/domain/services/accounts"
	handler "rollbringer/pkg/handlers/accounts"
	"rollbringer/pkg/repositories/nats"
	database "rollbringer/pkg/repositories/postgres/accounts"
	"rollbringer/pkg/repositories/spotify"
)

func init() {
	registeredFeatures["accounts"] = func() error {

		// Create PubSub repository.
		pubSubRepo, err := nats.NewPubSubRepository(config, logger.With("dependency", "nats-pubsub-repo"))
		if err != nil {
			return domain.Wrap(err, "cannot create PubSub repository", nil)
		}

		// Create Spotify repository.
		spotifyRepo := spotify.NewSpotifyRepository()

		// Create accounts database repository.
		accountsDBRepo, err := database.NewGamesDatabaseRepository(config, logger.With("dependency", "postgres-repo"), migrations)
		if err != nil {
			return domain.Wrap(err, "cannot create accounts database repository", nil)
		}

		svc := service.New(config, logger, pubSubRepo, accountsDBRepo, spotifyRepo)
		h := handler.New(config, logger.With("dependency", "http-api"), svc)

		features["accounts"] = feature{h, svc}
		return nil
	}
}

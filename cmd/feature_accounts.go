//go:build !noaccounts
// +build !noaccounts

package main

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth_spotify "golang.org/x/oauth2/spotify"

	"rollbringer/pkg/domain"
	service "rollbringer/pkg/domain/services/accounts"
	handler "rollbringer/pkg/handlers/accounts"
	database "rollbringer/pkg/repositories/database/accounts"
	"rollbringer/pkg/repositories/pubsub"
	"rollbringer/pkg/repositories/spotify"
)

func init() {
	registeredFeatures["accounts"] = func() error {

		oauthConfigGoogle := &oauth2.Config{
			Endpoint:     google.Endpoint,
			ClientID:     config.OauthGoogleClientID,
			ClientSecret: config.OauthGoogleClientSecret,
			Scopes:       []string{"openid", "profile", "email"},
		}

		oauthConfigSpotify := &oauth2.Config{
			Endpoint:     oauth_spotify.Endpoint,
			ClientID:     config.OauthSpotifyClientID,
			ClientSecret: config.OauthSpotifyClientSecret,
			Scopes:       []string{"user-read-private", "user-read-email"},
		}

		// Create PubSub repository.
		pubSubRepo, err := pubsub.NewPubSubRepository(config, logger.With("dependency", "nats-pubsub-repo"))
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

		svc := service.New(config, logger.With("dependency", "domain"), pubSubRepo, accountsDBRepo, spotifyRepo)
		h := handler.New(config, logger.With("dependency", "http-api"), oauthConfigGoogle, oauthConfigSpotify, svc)

		features["accounts"] = feature{h, svc}
		return nil
	}
}

package main

import (
	"net/http"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"rollbringer/pkg/api"
	"rollbringer/pkg/domain/service"
	"rollbringer/pkg/repositories/database"
	"rollbringer/pkg/repositories/pubsub"
)

type config struct {
	Address string `required:"true" split_words:"true"`

	PostgresAddress string `required:"true" split_words:"true"`

	RedisAddress  string `required:"true" split_words:"true"`
	RedisPassword string `required:"true" split_words:"true"`

	UserSessionTimeout time.Duration `required:"true" split_words:"true"`

	GoogleClientID     string `required:"true" split_words:"true"`
	GoogleClientSecret string `required:"true" split_words:"true"`

	RedirectURL string `required:"true" split_words:"true"`
}

func main() {

	// Create a logger.
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	logger := zerolog.New(os.Stdout)
	logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Create a configuration.
	var cfg config
	if err := envconfig.Process("APP", &cfg); err != nil {
		logger.Fatal().Stack().Err(err).Msg("Cannot generate configuration")
	}

	// Create a Database repository.
	db, err := database.New(cfg.PostgresAddress)
	if err != nil {
		logger.Fatal().Stack().Err(err).Msg("Cannot create database repository")
	}

	// Create a PubSub repository.
	ps, err := pubsub.New(&logger, cfg.RedisAddress, cfg.RedisPassword)
	if err != nil {
		logger.Fatal().Stack().Err(err).Msg("Cannot create pub-sub repository")
	}

	// Create an API.
	API := api.New(
		service.New(&logger, db, ps),
		&logger,
		&oauth2.Config{
			Endpoint:     google.Endpoint,
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       []string{"openid", "email"},
		},
	)

	logger.Info().Str("address", cfg.Address).Msg("Serving")

	// Start server.
	err = http.ListenAndServe(cfg.Address, API)
	logger.Fatal().Stack().Err(err).Msg("Server stopped")
}

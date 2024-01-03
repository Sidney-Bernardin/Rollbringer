package main

import (
	"context"
	"net/http"
	"os"
	"rollbringer/pkg/api"
	"rollbringer/pkg/repositories/database"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type config struct {
	Address   string `required:"true" split_words:"true"`
	DBAddress string `required:"true" split_words:"true"`

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

	// Create a sub-logger for the database.
	dbLoggerCtx := logger.With().
		Str("component", "database").
		Logger().WithContext(context.Background())

	// Create a database.
	db, err := database.New(cfg.DBAddress, dbLoggerCtx)
	if err != nil {
		logger.Fatal().Stack().Err(err).Msg("Cannot create database")
	}

	// Create an oauth2 configuration.
	googleOAuthConfig := &oauth2.Config{
		Endpoint:     google.Endpoint,
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes:       []string{"openid", "email"},
	}

	// Create and setup router with various services.
	router := createRouter(&api.API{
		DB:                db,
		Logger:            &logger,
		GoogleOAuthConfig: googleOAuthConfig,
	})

	logger.Info().Str("address", cfg.Address).Msg("Serving")

	// Start server.
	err = http.ListenAndServe(cfg.Address, router)
	logger.Fatal().Stack().Err(err).Msg("Cannot generate configuration")
}

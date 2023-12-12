package main

import (
	"embed"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"rollbringer/pkg/api"
	"rollbringer/pkg/database"
)

//go:embed templates
var templatesFS embed.FS

type Config struct {
	Address   string `required:"true" split_words:"true"`
	DBAddress string `required:"true" split_words:"true"`

	UserSessionTimeout time.Duration `required:"true" split_words:"true"`

	OauthClientID     string `required:"true" split_words:"true"`
	OauthClientSecret string `required:"true" split_words:"true"`
	OauthRedirectURL  string `required:"true" split_words:"true"`
}

func main() {

	// Create a logger.
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	logger := zerolog.New(os.Stdout)
	logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Create a configuration.
	var config Config
	if err := envconfig.Process("APP", &config); err != nil {
		logger.Fatal().Stack().Err(err).Msg("Cannot generate configuration")
	}

	// Create a database.
	db, err := database.New(config.DBAddress)
	if err != nil {
		logger.Fatal().Stack().Err(err).Msg("Cannot create database")
	}

	// Create an oauth2 configuration.
	googleOAuthConfig := &oauth2.Config{
		Endpoint:     google.Endpoint,
		ClientID:     config.OauthClientID,
		ClientSecret: config.OauthClientSecret,
		RedirectURL:  config.OauthRedirectURL,
		Scopes:       []string{"openid", "email"},
	}

	// Create the API.
	a, err := api.NewAPI(
		db,
		&logger,
		googleOAuthConfig,
		config.UserSessionTimeout,
		templatesFS, os.DirFS("static"))

	if err != nil {
		logger.Fatal().Stack().Err(err).Msg("Cannot create API")
	}

	var (
		errChan    = make(chan error)
		signalChan = make(chan os.Signal)
	)

	// Serve the API in another go-routine.
	logger.Info().Str("address", config.Address).Msg("Serving")
	go func() {
		errChan <- http.ListenAndServe(config.Address, a)
	}()

	// Have interrupt and termination signals sent to the signal-channel.
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-signalChan:
	case err := <-errChan:
		logger.Fatal().Stack().Err(err).Msg("API crashed")
	}
}

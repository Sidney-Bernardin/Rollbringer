package main

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"golang.org/x/net/websocket"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"rollbringer/pkg/domain/service"
	"rollbringer/pkg/handler"
	"rollbringer/pkg/repositories/database"
	"rollbringer/pkg/repositories/oauth"
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
	dbRepo, err := database.New(cfg.PostgresAddress)
	if err != nil {
		logger.Fatal().Stack().Err(err).Msg("Cannot create database repository")
	}

	// Create a PubSub repository.
	pubsubRepo, err := pubsub.New(cfg.RedisAddress, cfg.RedisPassword)
	if err != nil {
		logger.Fatal().Stack().Err(err).Msg("Cannot create pub-sub repository")
	}

	// Create an OAuth repository.
	oauthRepo := &oauth.OAuth{
		GoogleConfig: &oauth2.Config{
			Endpoint:     google.Endpoint,
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       []string{"openid", "profile", "email"},
		},
	}

	// Create a Service.
	svc := service.NewService(dbRepo, pubsubRepo, oauthRepo, &logger)

	// Create a Handler.
	h := &handler.Handler{
		Router:  chi.NewRouter(),
		Service: svc,
		Logger:  &logger,
	}

	h.Router.Use(h.Log, h.Instance)
	h.Router.Handle("/static/*", handleStaticDir())
	h.Router.Get("/", h.HandleHomePage)

	h.Router.With(h.AuthenticatePage).Get("/play", h.HandlePlayPage)
	h.Router.With(h.AuthenticatePage).Method("GET", "/ws", websocket.Handler(h.HandleWebSocket))

	h.Router.Route("/users", func(r chi.Router) {
		r.Get("/login", h.HandleLogin)
		r.Get("/consent-callback", h.HandleConsentCallback)
	})

	h.Router.Route("/games", func(r chi.Router) {
		r.Use(h.Authenticate)

		r.Post("/", h.HandleCreateGame)
		r.Delete("/{game_id}", h.HandleDeleteGame)
	})

	h.Router.Route("/play-materials", func(r chi.Router) {
		r.Use(h.Authenticate)

		r.Post("/pdfs", h.HandleCreatePDF)
		r.Get("/pdfs/{pdf_id}", h.HandleGetPDF)
		r.Delete("/pdfs/{pdf_id}", h.HandleDeletePDF)
	})

	// Start server.
	logger.Info().Str("address", cfg.Address).Msg("Serving")
	err = http.ListenAndServe(cfg.Address, h)
	logger.Fatal().Stack().Err(err).Msg("Server stopped")
}

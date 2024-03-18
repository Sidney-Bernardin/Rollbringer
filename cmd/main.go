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

	// Create a Service
	svc := &service.Service{
		DB:     db,
		PS:     ps,
		Logger: &logger,
	}

	// Create a Handler.
	h := &handler.Handler{
		Router:  chi.NewRouter(),
		Service: svc,
		Logger:  &logger,
		GoogleOAuthConfig: &oauth2.Config{
			Endpoint:     google.Endpoint,
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       []string{"openid", "email"},
		},
	}

	h.Router.Use(h.Log)
	h.Router.Handle("/static/*", handleStaticDir())

	h.Router.Get("/", h.HandleHomePage)
	h.Router.With(h.LightAuth).Get("/play", h.HandlePlayPage)
	h.Router.With(h.LightAuth).Method("GET", "/ws", websocket.Handler(h.HandleWebSocket))

	h.Router.Get("/users/login", h.HandleLogin)
	h.Router.Get("/users/consent-callback", h.HandleConsentCallback)

	h.Router.With(h.Auth).Post("/games", h.HandleCreateGame)
	h.Router.With(h.Auth).Get("/games", h.HandleGetGames)
	h.Router.With(h.Auth).Delete("/games/{game_id}", h.HandleDeleteGame)

	h.Router.With(h.Auth).Post("/play-materials/pdfs", h.HandleCreatePDF)
	h.Router.With(h.Auth).Get("/play-materials/pdfs", h.HandleGetPDFs)
	h.Router.With(h.Auth).Get("/play-materials/pdfs/{pdf_id}", h.HandleGetPDF)
	h.Router.With(h.Auth).Get("/play-materials/pdfs/{pdf_id}/{page_num}", h.HandleGetPDF)
	h.Router.With(h.Auth).Delete("/play-materials/pdfs/{pdf_id}", h.HandleDeletePDF)

	// Start server.
	logger.Info().Str("address", cfg.Address).Msg("Serving")
	err = http.ListenAndServe(cfg.Address, h)
	logger.Fatal().Stack().Err(err).Msg("Server stopped")
}

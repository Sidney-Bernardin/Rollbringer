package handler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/spotify"

	"rollbringer/pkg/domain"
	service "rollbringer/pkg/domain/services/accounts"
	"rollbringer/pkg/handlers"
)

type AccountsHandler struct {
	*handlers.Handler

	accountsSvc service.AccountsService

	oauthGoogleConfig  *oauth2.Config
	oauthSpotifyConfig *oauth2.Config
}

func New(config *domain.Config, logger *slog.Logger, gamesSvc service.AccountsService) http.Handler {
	h := &AccountsHandler{
		Handler: &handlers.Handler{
			Config:  config,
			Logger:  logger,
			Router:  chi.NewRouter(),
			Service: gamesSvc,
		},
		accountsSvc: gamesSvc,
		oauthGoogleConfig: &oauth2.Config{
			Endpoint:     google.Endpoint,
			ClientID:     config.OauthGoogleClientID,
			ClientSecret: config.OauthGoogleClientSecret,
			RedirectURL:  config.OauthGoogleRedirectURL,
			Scopes:       []string{"openid", "profile", "email"},
		},
		oauthSpotifyConfig: &oauth2.Config{
			Endpoint:     spotify.Endpoint,
			ClientID:     config.OauthSpotifyClientID,
			ClientSecret: config.OauthSpotifyClientSecret,
			RedirectURL:  config.OauthSpotifyRedirectURL,
			Scopes:       []string{"user-read-private", "user-read-email"},
		},
	}

	h.Router.Route("/login", func(r chi.Router) {
		r.Get("/google", h.handleOAuth(h.oauthGoogleConfig))
		r.With(h.mwOAuthCallback(h.oauthGoogleConfig)).
			Get("/google-callback", h.handleLoginCallbackGoogle)

		r.Get("/spotify", h.handleOAuth(h.oauthSpotifyConfig))
		r.With(h.mwOAuthCallback(h.oauthSpotifyConfig)).
			Get("/spotify-callback", h.handleLoginCallbackSpotify)
	})

	return h
}

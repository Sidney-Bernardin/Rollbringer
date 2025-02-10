package handler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"

	"rollbringer/pkg/domain"
	service "rollbringer/pkg/domain/services/accounts"
	"rollbringer/pkg/handlers"
)

type AccountsHandler struct {
	*handlers.Handler

	accountsSvc service.AccountsService
}

func New(config *domain.Config, logger *slog.Logger, oauthConfigGoogle, oauthConfigSpotify *oauth2.Config, gamesSvc service.AccountsService) http.Handler {
	h := &AccountsHandler{
		Handler: &handlers.Handler{
			Config:  config,
			Logger:  logger,
			Router:  chi.NewRouter(),
			Service: gamesSvc,
		},
		accountsSvc: gamesSvc,
	}

	h.Router.Use(h.MWLog)

	h.Router.Route("/signup", func(r chi.Router) {
		r.Route("/google", func(g chi.Router) {
			g.Use(h.mwOAuthConfig(oauthConfigGoogle, config.OauthGoogleSignupRedirectURL))

			g.Get("/", h.handleOAuth)
			g.With(h.mwOAuthCallback, h.mwCreateGoogleUser).Get("/callback", h.handleSignup)
		})

		r.Route("/spotify", func(s chi.Router) {
			s.Use(h.mwOAuthConfig(oauthConfigSpotify, config.OauthSpotifySignupRedirectURL))

			s.Get("/", h.handleOAuth)
			s.With(h.mwOAuthCallback, h.mwCreateSpotifyUser).Get("/callback", h.handleSignup)
		})
	})

	h.Router.Route("/signin", func(r chi.Router) {
		r.Route("/google", func(g chi.Router) {
			g.Use(h.mwOAuthConfig(oauthConfigGoogle, config.OauthGoogleSigninRedirectURL))

			g.Get("/", h.handleOAuth)
			g.With(h.mwOAuthCallback, h.mwCreateGoogleUser).Get("/callback", h.handleSignin)
		})

		r.Route("/spotify", func(s chi.Router) {
			s.Use(h.mwOAuthConfig(oauthConfigSpotify, config.OauthSpotifySigninRedirectURL))

			s.Get("/", h.handleOAuth)
			s.With(h.mwOAuthCallback, h.mwCreateSpotifyUser).Get("/callback", h.handleSignin)
		})
	})

	return h
}

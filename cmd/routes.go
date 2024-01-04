package main

import (
	"net/http"
	"os"
	"rollbringer/pkg/api"

	"github.com/go-chi/chi/v5"
	"golang.org/x/net/websocket"
)

func createRouter(a *api.API) chi.Router {

	router := chi.NewRouter()

	// Serve static files.
	router.Handle(
		"/static/*",
		http.StripPrefix("/static/", http.FileServer(http.FS(os.DirFS("static")))),
	)

	router.Route("/", func(r chi.Router) {
		r.Use(a.LightAuth)

		r.Get("/", a.HandleHomePage)
		r.Get("/game", a.HandleGamePage)
	})

	router.Route("/users", func(r chi.Router) {
		r.Get("/oauth-login", a.HandleOAuthLogin)
		r.Get("/oauth-consent-callback", a.HandleOAuthConsentCallback)
	})

	router.Route("/games", func(r chi.Router) {
		// r.Use(a.Auth)

		r.Method("GET", "/{id}/join", websocket.Handler(a.HandleJoinGame))
	})

	return router
}

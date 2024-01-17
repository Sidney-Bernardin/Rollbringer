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
	router.Use(a.Log)

	// Serve static files.
	router.Handle(
		"/static/*",
		http.StripPrefix("/static/", http.FileServer(http.FS(os.DirFS("static")))),
	)

	// Pages
	router.Route("/", func(r chi.Router) {
		r.Use(a.LightAuth)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/game", http.StatusTemporaryRedirect)
		})
		r.Get("/game", a.HandleGamePage)
	})

	// Users
	router.Route("/users", func(r chi.Router) {
		r.Get("/login", a.HandleLogin)
		r.Get("/consent-callback", a.HandleConsentCallback)
	})

	// Games
	router.Route("/games", func(r chi.Router) {
		r.With(a.Auth).Post("/", a.HandleCreateGame)
		r.With(a.Auth).Delete("/{game_id}", a.HandleDeleteGame)
		r.With(a.LightAuth).Method("GET", "/{game_id}/join", websocket.Handler(a.HandleJoinGame))
	})

	return router
}

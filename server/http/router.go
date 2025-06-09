package http

import (
	"embed"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed static
var static embed.FS

func (api *API) Router() http.Handler {
	r := chi.NewRouter()

	r.Use(api.mwLog)
	r.Handle("/static/*", http.FileServerFS(static))

	r.Post("/login", api.handleBasicLogin)
	r.Get("/login/google", api.handleGoogleLogin)
	r.Get("/login/google/callback", api.handleGoogleLoginCallback)
	r.Get("/logout", api.handleLogout)

	r.Route("/rooms", func(rr chi.Router) {
		authCSRF := rr.With(api.mwAuthenticate(true))
		authCSRF.Post("/", api.handleRoomsPost)
		authCSRF.Delete("/{room_id}", api.handleRoomsDelete)
	})

	auth := r.With(api.mwAuthenticate(false))
	auth.HandleFunc("/", api.handleHomePage)
	auth.HandleFunc("/play", api.handlePlayPage)

	return r
}

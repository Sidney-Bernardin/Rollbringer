package api

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"golang.org/x/net/websocket"
)

func (api *API) setRoutes() {

	api.router = chi.NewRouter()

	api.router.Use(api.Log)
	api.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/play", http.StatusTemporaryRedirect)
	})

	// Static files
	api.router.Handle(
		"/static/*",
		http.StripPrefix("/static/", http.FileServer(http.FS(os.DirFS("static")))),
	)

	api.router.Route("/play", func(r chi.Router) {
		r.Use(api.LightAuth)
		r.Get("/", api.HandlePlayPage)
		r.Method("GET", "/ws", websocket.Handler(api.HandlePlayWS))
	})

	api.router.Route("/users", func(r chi.Router) {
		r.Get("/login", api.HandleLogin)
		r.Get("/consent-callback", api.HandleConsentCallback)
	})

	api.router.Route("/games", func(r chi.Router) {
		r.With(api.Auth).Post("/", api.HandleCreateGame)
		r.With(api.Auth).Delete("/{game_id}", api.HandleDeleteGame)
	})
}

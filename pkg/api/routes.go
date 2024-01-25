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

	api.router.With(api.LightAuth).Get("/play", api.handlePlayPage)

	api.router.Route("/users", func(r chi.Router) {
		r.Get("/login", api.handleLogin)
		r.Get("/consent-callback", api.handleConsentCallback)
	})

	api.router.Route("/games", func(r chi.Router) {
		r.With(api.Auth).Post("/", api.handleCreateGame)
		r.With(api.Auth).Delete("/{game_id}", api.handleDeleteGame)
	})

	api.router.Route("/play-materials", func(r chi.Router) {
		r.Use(api.LightAuth)
		r.Method("GET", "/", websocket.Handler(api.handlePlayMaterials))

		// TODO: remove api.Auth when implementing guest-users.
		r.With(api.Auth).Post("/pdfs", api.handleCreatePDF)
		r.With(api.Auth).Delete("/pdfs/{pdf_id}", api.handleDeletePDF)
	})
}

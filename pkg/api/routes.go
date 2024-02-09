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
	api.router.Route("/", func(r chi.Router) {
		r.Use(api.LightAuth)
		r.Get("/play", api.handlePlayPage)
		r.Method("GET", "/ws", websocket.Handler(api.handleWebSocket))
	})

	api.router.Route("/users", func(r chi.Router) {
		r.Get("/login", api.handleLogin)
		r.Get("/consent-callback", api.handleConsentCallback)
	})

	api.router.Route("/games", func(r chi.Router) {
		r.Use(api.Auth)
		r.Post("/", api.handleCreateGame)
		r.Delete("/{game_id}", api.handleDeleteGame)
	})

	api.router.Route("/play-materials", func(r chi.Router) {
		r.Use(api.Auth)
		r.Post("/pdfs", api.handleCreatePDF)
		r.Get("/pdfs/{pdf_id}", api.handleGetPDF)
		r.Delete("/pdfs/{pdf_id}", api.handleDeletePDF)
	})
}

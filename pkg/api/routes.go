package api

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *api) doRoutes(staticFS fs.FS) {

	// Serve static files.
	a.router.Handle(
		"/static/*",
		http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))),
	)

	a.router.Route("/", func(r chi.Router) {
		r.Use(a.authenticate)

		r.Get("/", a.handleHomePage())
		r.Get("/game", a.handleGamePage())
	})

	a.router.Route("/users", func(r chi.Router) {
		r.Get("/oauth-login", a.handleOAuthLogin())
		r.Get("/oauth-consent-callback", a.handleOAuthConsentCallback())
	})

	a.router.Route("/games", func(r chi.Router) {
		r.Use(a.authenticate)

		r.Handle("/ws", a.handleWS())
	})
}

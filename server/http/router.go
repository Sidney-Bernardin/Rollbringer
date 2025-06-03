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
	r.Handle("/static/", http.FileServerFS(static))

	r.Get("/login/google", api.handleGoogleLogin)
	r.Get("/login/google/callback", api.handleGoogleLoginCallback)
	r.Get("/logout", api.handleLogout)

	public := r.With(api.mwAuthenticate(false))
	public.Get("/", api.handleHomePage)

	return r
}

package api

import (
	"html/template"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type api struct {
	router chi.Router
	tmpl   *template.Template
	logger *zerolog.Logger
}

func NewAPI(logger *zerolog.Logger, templatesFS, staticFS fs.FS) (a *api, err error) {

	// Create an api.
	a = &api{
		router: chi.NewRouter(),
		logger: logger,
	}

	// Parse templates.
	a.tmpl, err = template.ParseFS(templatesFS, "templates/*.html", "templates/pages/*.html")
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse templates")
	}

	a.doRoutes(staticFS)
	return a, errors.Wrap(err, "cannot parse templates")
}

func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

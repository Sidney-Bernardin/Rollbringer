package api

import (
	"html/template"
	"io/fs"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"

	"rollbringer/pkg/database"
)

const severErrorMsg = "Server Error"

var (
	errUnauthorized = errors.New("Unauthorized")
)

type api struct {
	router chi.Router
	tmpl   *template.Template

	db                *database.Database
	logger            *zerolog.Logger
	googleOAuthConfig *oauth2.Config

	userSessionTimeout time.Duration
}

func NewAPI(
	db *database.Database,
	logger *zerolog.Logger,
	googleOAuthConfig *oauth2.Config,
	userSessionTimeout time.Duration,
	templatesFS, staticFS fs.FS) (a *api, err error) {

	// Create an api.
	a = &api{
		router:            chi.NewRouter(),
		tmpl:              template.New("base").Funcs(templateUtils),
		db:                db,
		logger:            logger,
		googleOAuthConfig: googleOAuthConfig,
	}

	// Parse templates.
	a.tmpl, err = a.tmpl.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse templates")
	}

	a.doRoutes(staticFS)
	return a, errors.Wrap(err, "cannot parse templates")
}

func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func (a *api) executeTemplate(w http.ResponseWriter, name string, statusCode int, data any) {

	if errPgData, ok := data.(errorPageTmpl); ok {
		errPgData.StatusCode = statusCode
		errPgData.StatusText = http.StatusText(statusCode)
		data = errPgData

		if statusCode >= http.StatusInternalServerError {
			a.logger.Error().Stack().Err(errPgData.err).Msg(severErrorMsg)
		}
	}

	w.WriteHeader(statusCode)
	if err := a.tmpl.ExecuteTemplate(w, name, data); err != nil {
		a.logger.Error().Stack().Err(err).Msg(severErrorMsg)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

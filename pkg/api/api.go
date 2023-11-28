package api

import (
	"bytes"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/net/websocket"
	"golang.org/x/oauth2"

	"rollbringer/pkg/database"
)

type api struct {
	router chi.Router
	tmpl   *template.Template

	db          *database.Database
	logger      *zerolog.Logger
	oauthConfig *oauth2.Config
}

func NewAPI(
	db *database.Database,
	logger *zerolog.Logger,
	oauthConfig *oauth2.Config,
	templatesFS, staticFS fs.FS) (a *api, err error) {

	// Create an api.
	a = &api{
		router:      chi.NewRouter(),
		db:          db,
		logger:      logger,
		oauthConfig: oauthConfig,
	}

	// Parse templates.
	a.tmpl, err = template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse templates")
	}

	a.doRoutes(staticFS)
	return a, errors.Wrap(err, "cannot parse templates")
}

func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func (a *api) handleWS() websocket.Handler {

	type updatedPDFField struct {
		TextArea bool
		Type     string
		Name     string
		Value    string
	}

	return func(conn *websocket.Conn) {
		for {
			var msg string
			if err := websocket.Message.Receive(conn, &msg); err != nil {
				a.logger.Error().Stack().Err(err).Msg("Cannot recive ws msg")
				return
			}

			x := updatedPDFField{
				TextArea: false,
				Type:     "text",
				Name:     "testfield",
				Value:    "new value",
			}

			bbuf := bytes.Buffer{}
			if err := a.tmpl.ExecuteTemplate(&bbuf, "updated_pdf_fields", x); err != nil {
				a.logger.Error().Stack().Err(err).Msg("Cannot send ws template")
				return
			}

			conn.Write(bbuf.Bytes())
		}
	}
}

package api

import (
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"

	"rollbringer/pkg/domain/service"
)

const (
	wsStatusNormalClosure   = 1000
	wsStatusUnsupportedData = 1003
	wsStatusPolicyViolation = 1008
	wsStatusInternalError   = 1011
)

type API struct {
	router *chi.Mux
	Logger *zerolog.Logger
	GoogleOAuthConfig *oauth2.Config

	Service *service.Service
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.setRoutes()
	api.router.ServeHTTP(w, r)
}

func (api *API) render(w io.Writer, r *http.Request, component templ.Component, status int) {

	if rw, ok := w.(http.ResponseWriter); ok {
		rw.WriteHeader(status)
	}

	if err := component.Render(r.Context(), w); err != nil {
		err = errors.Wrap(err, "cannot render component")
		api.err(w, err, http.StatusInternalServerError, wsStatusInternalError)
	}
}

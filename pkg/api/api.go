package api

import (
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/net/websocket"
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
	router  *chi.Mux
	service *service.Service

	logger            *zerolog.Logger
	googleOAuthConfig *oauth2.Config
}

func New(service *service.Service, rootLogger *zerolog.Logger, googleOAuthConfig *oauth2.Config) *API {
	apiLogger := rootLogger.With().Str("component", "api").Logger()

	return &API{
		router:            chi.NewRouter(),
		service:           service,
		logger:            &apiLogger,
		googleOAuthConfig: googleOAuthConfig,
	}
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

func (api *API) closeConn(conn *websocket.Conn) {
	if err := conn.Close(); err != nil {
		api.logger.Error().Stack().Err(err).Msg("Cannot close connection")
	}
}

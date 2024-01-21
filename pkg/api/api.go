package api

import (
	"io"
	"net/http"
	"rollbringer/pkg/repositories/database"
	"rollbringer/pkg/repositories/pubsub"

	"github.com/a-h/templ"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
)

const (
	wsStatusNormalClosure   = 1000
	wsStatusUnsupportedData = 1003
	wsStatusPolicyViolation = 1008
	wsStatusInternalError   = 1011
)

type API struct {
	DB     *database.Database
	PubSub *pubsub.PubSub
	Logger *zerolog.Logger

	GoogleOAuthConfig *oauth2.Config
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

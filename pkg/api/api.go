package api

import (
	"context"
	"net/http"
	database "rollbringer/pkg/repositories/database"

	"github.com/a-h/templ"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/net/websocket"
	"golang.org/x/oauth2"
)

type API struct {
	DB     *database.Database
	Logger *zerolog.Logger

	GoogleOAuthConfig *oauth2.Config
}

// renderHTTP writes the component and status-code to the response-writer.
func (api *API) renderHTTP(w http.ResponseWriter, r *http.Request, component templ.Component, statusCode int) {
	w.WriteHeader(statusCode)
	if err := component.Render(r.Context(), w); err != nil {
		err = errors.Wrap(err, "cannot render component to HTTP response")
		api.err(w, err, statusCode)
	}
}

// renderWS writes the component to the WebSocket connection.
func (api *API) renderWS(ctx context.Context, conn *websocket.Conn, component templ.Component) {
	if err := component.Render(ctx, conn); err != nil {
		err = errors.Wrap(err, "cannot render component to HTTP response")
		api.Logger.Error().Stack().Err(err).Msg("Internal server error")
	}
}

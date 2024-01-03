package api

import (
	"errors"
	"net/http"
	database "rollbringer/pkg/repositories/database"
	"rollbringer/pkg/views"

	"github.com/a-h/templ"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
)

type API struct {
	DB     *database.Database
	Logger *zerolog.Logger

	GoogleOAuthConfig *oauth2.Config
}

func (api *API) render(w http.ResponseWriter, r *http.Request, component templ.Component, statusCode int) {
	w.WriteHeader(statusCode)
	if err := component.Render(r.Context(), w); err != nil {
		api.Logger.Error().Stack().Err(err).Msg("Server error")
	}
}

func (api *API) renderError(w http.ResponseWriter, r *http.Request, e error, statusCode int) {

	w.WriteHeader(statusCode)
	if statusCode >= 500 {
		api.Logger.Error().Stack().Err(e).Msg("Server error")
		e = errors.New("Internal Server Error")
	}

	if err := views.Error(e).Render(r.Context(), w); err != nil {
		api.Logger.Error().Stack().Err(err).Msg("Server error")
	}
}

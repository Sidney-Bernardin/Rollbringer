package handler

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

type Handler struct {
	Router  *chi.Mux
	Service *service.Service

	Logger            *zerolog.Logger
	GoogleOAuthConfig *oauth2.Config
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Router.ServeHTTP(w, r)
}

// render writes the Templ component to the io.Writer.
func (h *Handler) render(w io.Writer, r *http.Request, component templ.Component, httpStatus int) {

	if rw, ok := w.(http.ResponseWriter); ok {
		rw.WriteHeader(httpStatus)
	}

	if err := component.Render(r.Context(), w); err != nil {
		err = errors.Wrap(err, "cannot render component")
		h.err(w, err, http.StatusInternalServerError, wsStatusInternalError)
	}
}

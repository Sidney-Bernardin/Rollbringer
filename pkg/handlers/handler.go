package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"golang.org/x/net/websocket"

	"rollbringer/pkg/domain"
)

var statusCodes = map[domain.UserErrorType]int{
	domain.UsrErrTypeServerError:          http.StatusInternalServerError,
	domain.UsrErrTypeCannotProcessRequest: http.StatusUnprocessableEntity,
	domain.UsrErrTypeUnauthorized:         http.StatusUnauthorized,

	domain.UsrErrTypeGoogleUserDoesNotExists: http.StatusNotFound,
	domain.UsrErrTypeGoogleUserAlreadyExists: http.StatusConflict,

	domain.UsrErrTypeSpotifyUserDoesNotExists: http.StatusNotFound,
	domain.UsrErrTypeSpotifyUserAlreadyExists: http.StatusConflict,
}

type Handler struct {
	Config *domain.Config
	Logger *slog.Logger

	Router chi.Router

	Service domain.IService
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Router.ServeHTTP(w, r)
}

func (h *Handler) Err(w io.Writer, r *http.Request, err error) {
	userErr := domain.HandleError(r.Context(), h.Logger, slog.LevelError, err)

	switch w.(type) {
	case *websocket.Conn:
		h.Respond(w, r, 0, domain.Event{
			Operation: domain.OperationError,
			Payload:   userErr,
		})

	case http.ResponseWriter:
		h.Respond(w, r, statusCodes[userErr.Type], userErr)
	}
}

func (h *Handler) Respond(w io.Writer, r *http.Request, statusCode int, res any) {
	ctx, cancel := context.WithCancel(r.Context())
	*r = *r.WithContext(ctx)
	defer cancel()

	// If writing HTTP, write response header.
	if rw, ok := w.(http.ResponseWriter); ok {
		rw.WriteHeader(statusCode)
	}

	var err error
	switch res := res.(type) {

	// Templ response.
	case templ.Component:
		err = res.Render(r.Context(), w)
		err = domain.Wrap(err, "cannot render Templ response", nil)

	// JSON response.
	default:
		err = json.NewEncoder(w).Encode(res)
		err = domain.Wrap(err, "cannot JSON encode response", nil)
	}

	domain.HandleError(r.Context(), h.Logger, slog.LevelError, err)
}

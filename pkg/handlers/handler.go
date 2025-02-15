package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"

	"rollbringer/pkg/domain"
)

var statusCodes = map[domain.UserErrorType]int{
	domain.UsrErrTypeServerError:          http.StatusInternalServerError,
	domain.UsrErrTypeCannotProcessRequest: http.StatusUnprocessableEntity,
	domain.UsrErrTypeUnauthorized:         http.StatusUnauthorized,
	domain.UsrErrTypeRecordNotFound:       http.StatusNotFound,

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

func (h *Handler) State(r *http.Request) map[string]any {
	state, ok := r.Context().Value("state").(map[string]any)
	if !ok {
		state = map[string]any{}
		*r = *r.WithContext(context.WithValue(r.Context(), "state", state))
	}
	return state
}

func (h *Handler) authenticate(r *http.Request, checkCSRF bool) (*domain.Session, error) {
	var ctx = r.Context()

	cookie, err := r.Cookie("SESSION_ID")
	if err != nil {
		return nil, nil
	}

	sessionID, _ := uuid.Parse(cookie.Value)
	if sessionID == uuid.Nil {
		return nil, nil
	}

	session, err := h.Service.GetSession(ctx, sessionID)
	if err != nil {
		if userErr, ok := errors.Cause(err).(*domain.UserError); ok && userErr.Type == domain.UsrErrTypeRecordNotFound {
			return nil, nil
		}

		return nil, domain.Wrap(err, "cannot get session", nil)
	}

	if checkCSRF && r.Header.Get("CSRF-TOKEN") != session.CSRFToken {
		return nil, nil
	}

	return session, nil
}

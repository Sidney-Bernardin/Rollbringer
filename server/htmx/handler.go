package htmx

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"rollbringer/server"
	"rollbringer/server/domain"
	"rollbringer/server/domain/play"
	"rollbringer/server/domain/play/results"
	"rollbringer/server/htmx/views"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"golang.org/x/net/websocket"
)

var errCodes = map[domain.DomainErrorType]int{
	results.DomainErrorTypeUUIDInvalid:     http.StatusBadRequest,
	results.DomainErrorTypeRoomNameInvalid: http.StatusBadRequest,
	results.DomainErrorTypeRoomNameTaken:   http.StatusConflict,
}

type handler struct {
	router chi.Router

	playSvc play.Service
	playDB  play.Database
	playBkr play.Broker

	logger *slog.Logger
	config *server.Config
}

func NewHandler(logger *slog.Logger) *handler {
	h := &handler{
		router: chi.NewRouter(),
		logger: logger,
	}

	h.router.With(h.mwInstance("home-page")).
		Get("/pages", h.handlePageHome)

	return h
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *handler) err(w io.Writer, r *http.Request, err error) {
	var ctx = r.Context()

	var domainErr *domain.DomainError
	if !errors.As(err, domainErr) {
		h.logServerError(ctx, err)
		h.respond(w, r, http.StatusInternalServerError, &views.ProblemDetail{
			Instance: ctx.Value("instance").(string),
			Type:     "internal_server_error",
		})
		return
	}

	problemDetail := &views.ProblemDetail{
		Instance: ctx.Value("instance").(string),
		Type:     string(domainErr.Type),
		Detail:   domainErr.Description,
	}

	switch w.(type) {
	case *websocket.Conn:
		h.respond(w, r, 0, problemDetail)

	case http.ResponseWriter:
		h.respond(w, r, errCodes[domainErr.Type], problemDetail)
	}
}

func (h *handler) respond(w io.Writer, r *http.Request, statusCode int, res any) {

	var (
		state = h.state(r)
		ctx   = r.Context()
	)

	if rw, ok := w.(http.ResponseWriter); ok {
		rw.WriteHeader(statusCode)
		state["status_code"] = statusCode
	}

	var err error
	switch res := res.(type) {

	case templ.Component:
		err = res.Render(r.Context(), w)
		err = errors.Join(err, errors.New("cannot render Templ response"))

	default:
		err = json.NewEncoder(w).Encode(res)
		err = errors.Join(err, errors.New("cannot JSON encode response"))
	}

	h.logServerError(ctx, err)
}

func (h *handler) state(r *http.Request) map[string]any {
	state, ok := r.Context().Value("state").(map[string]any)
	if !ok {
		state = map[string]any{}
		*r = *r.WithContext(context.WithValue(r.Context(), "state", state))
	}
	return state
}

func (h *handler) logServerError(ctx context.Context, err error) {
	h.logger.Log(ctx, server.LevelError,
		"Internal Server Error", "err", err.Error())
}

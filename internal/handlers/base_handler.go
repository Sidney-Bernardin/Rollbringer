package handlers

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"rollbringer/internal"
	"rollbringer/internal/config"
	"rollbringer/internal/services"
)

type BaseHandler struct {
	Config *config.Config
	Logger *slog.Logger

	Router *chi.Mux

	BaseService services.BaseServicer
}

func (h *BaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Router.ServeHTTP(w, r)
}

func (h *BaseHandler) authenticate(r *http.Request, sessionView internal.SessionView, checkCSRFToken bool) (*internal.Session, error) {
	var ctx = r.Context()

	cookie, err := r.Cookie("SESSION_ID")
	if err != nil {
		return nil, internal.NewProblemDetail(ctx, internal.PDOpts{Type: internal.PDTypeUnauthorized})
	}

	sessionID, err := uuid.Parse(cookie.Value)
	if err != nil {
		return nil, internal.NewProblemDetail(ctx, internal.PDOpts{Type: internal.PDTypeUnauthorized})
	}

	session, err := h.BaseService.Authenticate(ctx, sessionID, sessionView, checkCSRFToken, r.Header.Get("CSRF-Token"))
	return session, errors.Wrap(err, "cannot authenticate")
}

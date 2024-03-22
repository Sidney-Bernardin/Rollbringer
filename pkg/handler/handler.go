package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"

	"rollbringer/pkg/domain"
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

func (h *Handler) processServerError(err error) {
	if _, ok := errors.Cause(err).(*domain.ProblemDetail); ok {
		return
	}

	h.Logger.Error().Stack().Err(err).Msg("Server error")

	dest := &err

	if event, ok := err.(*domain.EventError); ok {
		dest = &event.Err
	}

	*dest = &domain.ProblemDetail{
		Type: domain.PDTypeServerError,
	}
}

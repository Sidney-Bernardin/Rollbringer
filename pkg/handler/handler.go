package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"rollbringer/pkg/domain/service"
)

type Handler struct {
	Router  *chi.Mux
	Service *service.Service

	Logger *zerolog.Logger
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Router.ServeHTTP(w, r)
}

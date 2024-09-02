package handlers

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"rollbringer/internal/config"
	"rollbringer/internal/services"
)

type Handler struct {
	Config *config.Config
	Logger *slog.Logger

	Router *chi.Mux

	svc services.Servicer
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Router.ServeHTTP(w, r)
}

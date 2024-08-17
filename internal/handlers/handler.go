package handlers

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"rollbringer/internal/config"
)

type Handler struct {
	Logger *slog.Logger
	Config *config.Config

	Router *chi.Mux
}

func NewHandler(cfg *config.Config, logger *slog.Logger) *Handler {
	return &Handler{
		Logger: logger,
		Config: cfg,
		Router: chi.NewRouter(),
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Router.ServeHTTP(w, r)
}

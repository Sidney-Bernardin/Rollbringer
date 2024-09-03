package games

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"rollbringer/internal/handlers"
	"rollbringer/internal/config"
	service "rollbringer/internal/services/games"
)

type handler struct {
	*handlers.Handler

	svc service.Service
}

func NewHandler(cfg *config.Config, logger *slog.Logger, svc service.Service) http.Handler {
	h := &handler{
		Handler: &handlers.Handler{
			Config: cfg,
			Logger: logger,
			Router: chi.NewRouter(),
		},
		svc: svc,
	}

	h.Router.Use(h.Log, h.Instance, h.Authenticate)
	h.Router.Post("/", h.HandleCreateGame)
	h.Router.Delete("/{game_id}", h.HandleDeleteGame)

	return h
}

func (h *handler) HandleCreateGame(w http.ResponseWriter, r *http.Request) {}
func (h *handler) HandleDeleteGame(w http.ResponseWriter, r *http.Request) {}

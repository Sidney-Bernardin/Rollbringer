package handler

import (
	"log/slog"
	"net/http"

	"rollbringer/internal/config"
	"rollbringer/internal/handlers"
	service "rollbringer/internal/services/games"
)

type gamesHandler struct {
	*handlers.Handler

	svc service.GamesService
}

func New(cfg *config.Config, logger *slog.Logger, svc service.GamesService) *gamesHandler {
	h := &gamesHandler{
		Handler: handlers.NewHandler(cfg, logger),
		svc:     svc,
	}

	h.Router.Use(h.Log, h.Instance, h.Authenticate)
	h.Router.Post("/", h.HandleCreateGame)
	h.Router.Delete("/{game_id}", h.HandleDeleteGame)

	return h
}

func (h *gamesHandler) HandleCreateGame(w http.ResponseWriter, r *http.Request) {}
func (h *gamesHandler) HandleDeleteGame(w http.ResponseWriter, r *http.Request) {}

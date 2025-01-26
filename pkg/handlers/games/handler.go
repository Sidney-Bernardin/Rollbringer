package handler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"rollbringer/pkg/domain"
	service "rollbringer/pkg/domain/services/games"
	"rollbringer/pkg/handlers"
)

type GamesHandler struct {
	*handlers.Handler

	gamesSvc service.GamesService
}

func New(config *domain.Config, logger *slog.Logger, gamesSvc service.GamesService) http.Handler {
	h := &GamesHandler{
		Handler: &handlers.Handler{
			Config:  config,
			Logger:  logger,
			Router:  chi.NewRouter(),
			Service: gamesSvc,
		},
		gamesSvc: gamesSvc,
	}

	return h
}

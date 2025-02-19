package handler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"rollbringer/pkg/domain"
	service "rollbringer/pkg/domain/services/play"
	"rollbringer/pkg/handlers"
)

type playHandler struct {
	*handlers.Handler

	playSvc service.PlayService
}

func New(config *domain.Config, logger *slog.Logger, playSvc service.PlayService) http.Handler {
	h := &playHandler{
		Handler: &handlers.Handler{
			Config:  config,
			Logger:  logger,
			Router:  chi.NewRouter(),
			Service: playSvc,
		},
		playSvc: playSvc,
	}

	h.Router.Use(h.MWLog)

	h.Router.Route("/rooms", func(r chi.Router) {
		auth := r.With(h.MWAuthenticate(true, true, false))

		auth.Post("/", h.handleRoomsPost)
		auth.Delete("/{room_id}", h.handleRoomsDelete)
	})

	h.Router.Route("/boards", func(r chi.Router) {
		auth := r.With(h.MWAuthenticate(true, true, false))

		auth.Post("/", h.handleBoardsPost)
		r.Get("/{board_id}", h.handleBoardGet)
	})

	return h
}

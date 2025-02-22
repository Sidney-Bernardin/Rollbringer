package handler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/net/websocket"

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
		csrf := r.With(h.MWAuthenticate(true, true, false))
		noCsrf := r.With(h.MWAuthenticate(true, false, false))

		csrf.Post("/", h.handleRoomsPost)
		csrf.Delete("/{room_id}", h.handleRoomsDelete)
		noCsrf.Handle("/{room_id}/ws", websocket.Handler(h.handleRoomsWebSocket))
	})

	h.Router.Route("/boards", func(r chi.Router) {
		csrf := r.With(h.MWAuthenticate(true, true, false))

		csrf.Post("/", h.handleBoardsPost)
		r.Get("/{board_id}", h.handleBoardGet)
	})

	return h
}

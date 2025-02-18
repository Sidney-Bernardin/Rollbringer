package handler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

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

	return h
}

func (h *playHandler) handleRoomsPost(w http.ResponseWriter, r *http.Request) {

	var (
		state = h.State(r)
		ctx   = r.Context()

		room = &domain.Room{
			Name: r.FormValue("name"),
		}
	)

	if err := h.playSvc.CreateRoom(ctx, state["session"].(*domain.Session), room); err != nil {
		h.Err(w, r, domain.Wrap(err, "cannot create room", nil))
		return
	}

	h.Respond(w, r, http.StatusOK, room.ID)
}

func (h *playHandler) handleRoomsDelete(w http.ResponseWriter, r *http.Request) {

	var (
		state = h.State(r)
		ctx   = r.Context()

		roomID, _ = uuid.Parse(chi.URLParam(r, "room_id"))
	)

	if err := h.playSvc.DeleteRoom(ctx, state["session"].(*domain.Session), roomID); err != nil {
		h.Err(w, r, domain.Wrap(err, "cannot delete room", nil))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

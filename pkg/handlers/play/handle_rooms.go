package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/handlers/play/views"
)

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

	h.Respond(w, r, http.StatusOK, views.RoomListItem(room))
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

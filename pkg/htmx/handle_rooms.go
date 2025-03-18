package htmx

import (
	"errors"
	"net/http"
	"rollbringer/pkg/domain/play/commands"
	"rollbringer/pkg/domain/play/results"

	"github.com/go-chi/chi/v5"
)

func (h *handler) RoomGet(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	// Get the room-ID.
	roomID, err := commands.NewUUID(chi.URLParam(r, "room_id"))
	if err != nil {
		h.err(w, r, errors.Join(err, errors.New("cannot create room-id")))
		return
	}

	// Initialize the result.
	var result any
	switch r.URL.Query().Get("v") {
	case "info":
		result = results.RoomInfo{}
	case "list-item":
		result = results.RoomListItem{}
	}

	// Get the room by ID.
	if err := h.playDB.RoomGetByID(ctx, roomID, &result); err != nil {
		h.err(w, r, errors.Join(err, errors.New("cannot get room by id")))
		return
	}

	h.respond(w, r, http.StatusOK, result)
}


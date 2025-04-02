package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"rollbringer/src/domain/play"
)

func (svr *server) handleRoomGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var room any
		switch r.URL.Query().Get("v") {
		case "info":
			room = play.RoomInfo{}
		case "list-item":
			room = play.RoomListItem{}
		}

		if err := svr.play.RoomGetByID(r.Context(), &room, chi.URLParam(r, "room_id")); err != nil {
			svr.err(w, r, errors.Wrap(err, "cannot get room by ID"))
			return
		}

		svr.respond(w, r, http.StatusOK, room)
	})
}

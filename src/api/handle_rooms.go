package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"rollbringer/src/domain/play"
)

func (svr *server) handleRoomCreate(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	// Initialize the room.
	var room any
	switch r.URL.Query().Get("v") {
	case "info":
		room = play.RoomInfo{}
	case "list-item":
		room = play.RoomListItem{}
	}

	err := svr.play.RoomCreate(ctx, room, &play.ArgsRoomCreate{
		OwnerID: r.FormValue("owner_id"),
		Name:    r.FormValue("name"),
	})

	if err != nil {
		svr.err(w, r, errors.Wrap(err, "cannot create room"))
		return
	}

	svr.respond(w, r, http.StatusOK, room)
}

func (svr *server) handleRoomGet(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	// Initialize the room.
	var room any
	switch r.URL.Query().Get("v") {
	case "info":
		room = play.RoomInfo{}
	case "list-item":
		room = play.RoomListItem{}
	}

	if err := svr.play.RoomGetByID(ctx, &room, chi.URLParam(r, "room_id")); err != nil {
		svr.err(w, r, errors.Wrap(err, "cannot get room by ID"))
		return
	}

	svr.respond(w, r, http.StatusOK, room)
}

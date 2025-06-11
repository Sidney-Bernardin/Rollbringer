package http

import (
	"net/http"

	"github.com/Sidney-Bernardin/Rollbringer/server"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/cache"
	"github.com/Sidney-Bernardin/Rollbringer/web/components"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

func (api *API) handleRoomsPost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, _ := ctx.Value("session").(*cache.Session)
	if session == nil {
		api.err(w, r, &server.UserError{Type: server.UserErrorTypeUnauthorized})
		return
	}

	room, err := api.Service.CreateRoom(r.Context(), session.UserID, r.FormValue("name"))
	if err != nil {
		api.err(w, r, errors.Wrap(err, "cannot create room"))
		return
	}

	api.respond(w, r, http.StatusOK, components.RoomCard(room.ID, room.Name, room.UserIds, room.Usernames, room.ProfilePictures, room.Permisions))
}

func (api *API) handleRoomsDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, _ := ctx.Value("session").(*cache.Session)
	if session == nil {
		api.err(w, r, &server.UserError{Type: server.UserErrorTypeUnauthorized})
		return
	}

	roomID, err := server.ParseUUID(chi.URLParam(r, "room_id"))
	if err != nil {
		api.err(w, r, errors.Wrap(err, "cannot parse room-ID"))
		return
	}

	if err := api.Service.DeleteRoom(r.Context(), session.UserID, roomID); err != nil {
		api.err(w, r, errors.Wrap(err, "cannot delete room"))
		return
	}
}

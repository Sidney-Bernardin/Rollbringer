package api

import (
	"net/http"

	"github.com/pkg/errors"

	"rollbringer/src"
	"rollbringer/src/api/views"
	accountModels "rollbringer/src/services/accounts/models"
	playModels "rollbringer/src/services/play/models"
)

func (svr *server) handleRoomCreate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var session, _ = r.Context().Value("session_info").(*accountModels.Session)

		roomName, err := playModels.ParseRoomName(r.FormValue("name"))
		if err != nil {
			svr.err(w, r, errors.Wrap(err, "cannot parse room-name"))
			return
		}

		room := &playModels.Room{
			ID:   src.NewUUID(),
			Name: roomName,
			Users: map[src.UUID]*playModels.RoomUser{
				session.UserID: {
					UserID:     session.UserID,
					Permisions: []playModels.RoomUserPermision{playModels.RoomUserPermisionOwner, playModels.RoomUserPermisionGameMaster},
				},
			},
		}

		if err = svr.play.CreateRoom(r.Context(), room); err != nil {
			svr.err(w, r, errors.Wrap(err, "cannot create room"))
			return
		}

		svr.respond(w, r, http.StatusOK, views.RoomCard(room))
	})
}

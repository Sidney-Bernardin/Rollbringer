package api

import (
	"net/http"
	"rollbringer/src/api/views"
	account_models "rollbringer/src/services/accounts/models"
	play_models "rollbringer/src/services/play/models"

	"github.com/pkg/errors"
)

func (svr *server) handleRoomCreate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var session, _ = r.Context().Value("session").(*account_models.Session)

		room, err := play_models.NewRoom(session.UserID, r.FormValue("name"))
		if err != nil {
			svr.err(w, r, errors.Wrap(err, "cannot create room"))
			return
		}

		if err = svr.playDatabase.CreateRoom(r.Context(), room); err != nil {
			svr.err(w, r, errors.Wrap(err, "cannot create room"))
			return
		}

		svr.respond(w, r, http.StatusOK, views.RoomCard(room, []*account_models.User{session.User}))
	})
}

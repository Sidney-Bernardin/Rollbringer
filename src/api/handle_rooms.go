package api

import (
	"net/http"
)

func (svr *server) handleRoomCreate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// var room any
		// switch r.URL.Query().Get("v") {
		// case "info":
		// 	room = play.RoomInfo{}
		// case "list-item":
		// 	room = play.RoomListItem{}
		// }
		//
		// if err := svr.play.RoomCreate(r.Context(), &room, nil); err != nil {
		// 	svr.err(w, r, errors.Wrap(err, "cannot create room"))
		// 	return
		// }

		svr.respond(w, r, http.StatusOK, r.Context().Value("session_info"))
	})
}

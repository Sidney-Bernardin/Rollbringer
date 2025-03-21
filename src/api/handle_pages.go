package api

import (
	"net/http"

	"github.com/pkg/errors"

	"rollbringer/src/api/views/pages"
	"rollbringer/src/domain/play"
)

func (svr *server) handlePageHome(w http.ResponseWriter, r *http.Request) {

	var (
		ctx  = r.Context()
		page = &pages.HomeData{
			RoomListItem: &play.RoomListItem{},
		}
	)

	// Get the room by ID.
	if err := svr.play.RoomGetByID(ctx, page.RoomListItem, r.URL.Query().Get("r")); err != nil {
		svr.err(w, r, errors.Wrap(err, "cannot get room by ID"))
		return
	}

	svr.respond(w, r, http.StatusOK, pages.Home(page))
}

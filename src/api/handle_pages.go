package api

import (
	"net/http"

	"github.com/pkg/errors"

	"rollbringer/src/api/views/pages"
	"rollbringer/src/domain/accounts"
	"rollbringer/src/domain/play"
)

func (svr *server) handlePageHome(w http.ResponseWriter, r *http.Request) {

	var (
		ctx  = r.Context()
		page = &pages.HomeData{
			UserInfo:     &accounts.UserInfo{},
			RoomListItem: &play.RoomListItem{},
		}
	)

	if err := svr.accounts.UserGetByUsername(ctx, page.UserInfo, r.URL.Query().Get("u")); err != nil {
		svr.err(w, r, errors.Wrap(err, "cannot get user by username"))
		return
	}

	if err := svr.play.RoomGetByID(ctx, page.RoomListItem, r.URL.Query().Get("r")); err != nil {
		svr.err(w, r, errors.Wrap(err, "cannot get room by ID"))
		return
	}

	svr.respond(w, r, http.StatusOK, pages.Home(page))
}

package api

import (
	"net/http"

	"github.com/pkg/errors"

	"rollbringer/src"
	"rollbringer/src/api/views/pages"
	accountModels "rollbringer/src/services/accounts/models"
	playModels "rollbringer/src/services/play/models"
)

func (svr *server) handlePageHome() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			ctx        = r.Context()
			session, _ = ctx.Value("session").(*accountModels.Session)
			err        error
			page       = &pages.HomeData{
				Session:   session,
				Rooms:     []*playModels.Room{},
				RoomUsers: map[src.UUID][]*accountModels.User{},
			}
		)

		if session == nil {
			svr.respond(w, r, http.StatusOK, pages.Home(page))
			return
		}

		// Get the user's rooms.
		page.Rooms, err = svr.playDB.GetRoomsByUserID(ctx, session.UserID)
		if err != nil {
			svr.err(w, r, errors.Wrap(err, "cannot get rooms by user-ID"))
			return
		}

		roomIDs := make([]src.UUID, 0, len(page.Rooms))
		for _, room := range page.Rooms {
			roomIDs = append(roomIDs, room.ID)
		}

		// Get users for each room.
		page.RoomUsers, err = svr.accountsDB.GetUsersByRoomIDs(ctx, roomIDs...)
		if err != nil {
			svr.err(w, r, errors.Wrap(err, "cannot get users by room-IDs"))
			return
		}

		svr.respond(w, r, http.StatusOK, pages.Home(page))
	})
}

func (svr *server) handlePagePlay() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var session, _ = r.Context().Value("session").(*accountModels.Session)

		svr.respond(w, r, http.StatusOK, pages.Play(&pages.PlayData{
			Session: session,
		}))
	})
}

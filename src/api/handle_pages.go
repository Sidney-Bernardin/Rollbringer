package api

import (
	"net/http"

	"rollbringer/src/api/views/pages"
	"rollbringer/src/services/accounts/models"
)

func (svr *server) handlePageHome() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var session, _ = r.Context().Value("session").(*models.Session)

		svr.respond(w, r, http.StatusOK, pages.Home(&pages.HomeData{
			Session: session,
		}))
	})
}

func (svr *server) handlePagePlay() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var session, _ = r.Context().Value("session").(*models.Session)

		svr.respond(w, r, http.StatusOK, pages.Play(&pages.PlayData{
			Session: session,
		}))
	})
}

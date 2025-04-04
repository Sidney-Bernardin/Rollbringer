package api

import (
	"net/http"

	"rollbringer/src/api/views/pages"
	"rollbringer/src/domain/accounts"
)

func (svr *server) handlePageHome() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sessionInfo, _ = r.Context().Value("session_info").(*accounts.ViewSessionInfo)

		svr.respond(w, r, http.StatusOK, pages.Home(&pages.HomeData{
			SessionInfo: sessionInfo,
		}))
	})
}

func (svr *server) handlePagePlay() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sessionInfo, _ = r.Context().Value("session_info").(*accounts.ViewSessionInfo)

		svr.respond(w, r, http.StatusOK, pages.Play(&pages.PlayData{
			SessionInfo: sessionInfo,
		}))
	})
}

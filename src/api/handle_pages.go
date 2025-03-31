package api

import (
	"net/http"

	"rollbringer/src/api/views/pages"
)

func (svr *server) handlePageHome() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		svr.respond(w, r, http.StatusOK, pages.Home(&pages.HomeData{}))
	})
}

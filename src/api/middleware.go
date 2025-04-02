package api

import (
	"net/http"

	"rollbringer/src"
)

type middleware func(http.Handler) http.Handler

func mw(mm ...middleware) middleware {
	return func(next http.Handler) http.Handler {
		for _, m := range mm {
			next = m(next)
		}
		return next
	}
}

/////

func (svr *server) mwLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		svr.log.Log(r.Context(), src.LevelInfo, "New request",
			"method", r.Method,
			"path", r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

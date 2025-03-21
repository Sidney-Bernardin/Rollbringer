package api

import (
	"net/http"
	"rollbringer/src"
)

func (svr *server) mwLog() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			svr.log.Log(r.Context(), src.LevelInfo, "New request",
				"method", r.Method,
				"path", r.URL.Path)

			next.ServeHTTP(w, r)
		})
	}
}

func (svr *server) mwInstance(instance string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			svr.state(r)["instance"] = instance
			next.ServeHTTP(w, r)
		})
	}
}

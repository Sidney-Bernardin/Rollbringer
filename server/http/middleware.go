package http

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Sidney-Bernardin/Rollbringer/server"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/cache"

	"github.com/pkg/errors"
)

func (api *API) mwLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer next.ServeHTTP(w, r)
		api.Log.Log(r.Context(), slog.LevelInfo, "New request",
			"method", r.Method,
			"path", r.URL.Path)
	})
}

func (api *API) mwAuthenticate(csrf bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			cSessionID, err := r.Cookie("SESSION_ID")
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			sessionID, err := server.ParseUUID(cSessionID.Value)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			session, err := api.Service.Cache.GetSession(ctx, sessionID)
			if err != nil {
				if errors.Is(err, cache.ErrNotFound) {
					next.ServeHTTP(w, r)
				} else {
					api.err(w, r, errors.Wrap(err, "cannot get session"))
				}
				return
			}

			if csrf && session.CSRFToken != r.Header.Get("CSRF-Token") {
				return
			}

			ctx = context.WithValue(ctx, "session", session)
			*r = *r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

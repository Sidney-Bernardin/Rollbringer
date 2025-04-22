package api

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	"rollbringer/src"
	"rollbringer/src/services/accounts/models"
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

func (svr *server) mwLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		svr.log.Log(r.Context(), src.LevelInfo, "New request",
			"method", r.Method,
			"path", r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

func (svr *server) mwAuth(required, checkCSRF bool, redirectURL string) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()

			var csrfToken *string
			if checkCSRF {
				h := r.Header.Get("CSRF-Token")
				csrfToken = &h
			}

			var session *models.Session
			if sessionID, err := r.Cookie("SESSION_ID"); err == nil {
				if session, err = svr.accounts.Auth(ctx, sessionID.Value, csrfToken); err != nil {
					svr.err(w, r, errors.Wrap(err, "cannot authenticate"))
					return
				}
			}

			if session == nil && required {
				if redirectURL != "" {
					http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
				} else {
					svr.err(w, r, &src.ExternalError{Type: src.ExternalErrorTypeUnauthorized})
				}
				return
			}

			*r = *r.WithContext(context.WithValue(ctx, "session", session))
			next.ServeHTTP(w, r)
		})
	}
}

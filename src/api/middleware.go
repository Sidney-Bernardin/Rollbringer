package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"rollbringer/src"
	"rollbringer/src/domain"
	"rollbringer/src/domain/services/accounts"
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

			var (
				ctx     = r.Context()
				session *accounts.Session
			)

			defer func() {
				if session == nil && required {
					if redirectURL != "" {
						http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
					} else {
						svr.err(w, r, &domain.ExternalError{Type: domain.ExternalErrorTypeUnauthorized})
					}
					return
				}

				*r = *r.WithContext(context.WithValue(ctx, "session", session))
				next.ServeHTTP(w, r)
			}()

			// Get the session-ID.
			cSessionID, err := r.Cookie("SESSION_ID")
			if err != nil {
				return
			}

			// Parse the session-ID.
			sessionID, err := uuid.Parse(cSessionID.Value)
			if err != nil {
				return
			}

			// Get the session.
			if csrfToken := r.Header.Get("CSRF-Token"); checkCSRF {
				session, err = svr.accountsDatabase.GetSessionBySessionIDAndCSRFToken(ctx, sessionID, csrfToken)
			} else {
				session, err = svr.accountsDatabase.GetSessionBySessionID(ctx, sessionID)
			}

			if err != nil && !errors.Is(err, domain.ErrEntityNotFound) {
				svr.err(w, r, errors.Wrap(err, "cannot authenticate"))
				return
			}
		})
	}
}

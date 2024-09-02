package pages

import (
	"context"
	"net/http"
	"rollbringer/internal"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (h *handler) AuthenticatePage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("SESSION_ID")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/users/login", http.StatusTemporaryRedirect)
				return
			}

			h.Err(w, r, errors.Wrap(err, "cannot get SESSION_ID cookie"))
			return
		}
		sessionID, _ := uuid.Parse(cookie.Value)

		session, err := h.svc.GetSession(r.Context(), sessionID)
		if err != nil {
			if internal.IsDetailed(err, internal.PDTypeUnauthorized) {
				http.Redirect(w, r, "/users/login", http.StatusTemporaryRedirect)
				return
			}

			h.Err(w, r, errors.Wrap(err, "cannot get session"))
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "session", session))
		next.ServeHTTP(w, r)
	})
}

package handler

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"rollbringer/internal"
)

func (h *gamesHandler) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		stCookie, err := r.Cookie("SESSION_ID")
		if err != nil {
			if err == http.ErrNoCookie {
				h.Err(w, r, &internal.ProblemDetail{
					Type: internal.PDTypeUnauthorized,
				})
				return
			}

			h.Err(w, r, errors.Wrap(err, "cannot get SESSION_ID cookie"))
			return
		}
		sessionID, _ := uuid.Parse(stCookie.Value)

		session, err := h.svc.Authenticate(r.Context(), sessionID, r.Header.Get("CSRF-Token"))
		if err != nil {
			if internal.IsDetailed(err, internal.PDTypeUnauthorized) {
				h.Err(w, r, &internal.ProblemDetail{
					Type: internal.PDTypeUnauthorized,
				})
				return
			}

			h.Err(w, r, errors.Wrap(err, "cannot get session"))
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "session", session))
		next.ServeHTTP(w, r)
	})
}

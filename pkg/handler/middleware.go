package handler

import (
	"context"
	"net/http"
	"rollbringer/pkg/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (h *Handler) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Log the request.
		h.Logger.Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Msg("New request")

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get the session-ID from the cookie.
		stCookie, err := r.Cookie("SESSION_ID")
		if err != nil {
			if err == http.ErrNoCookie {
				next.ServeHTTP(w, r)
				return
			}

			h.renderErr(w, r, http.StatusInternalServerError, errors.Wrap(err, "cannot get CSRF_Token cookie"))
			return
		}
		sessionID, _ := uuid.Parse(stCookie.Value)

		// Get the session.
		session, err := h.Service.GetSession(r.Context(), sessionID, nil)
		if err != nil && domain.IsProblemDetail(err, domain.PDTypeUnauthorized) {
			h.err(w, r, errors.Wrap(err, "cannot get session"))
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "session", session))
		next.ServeHTTP(w, r)
	})
}

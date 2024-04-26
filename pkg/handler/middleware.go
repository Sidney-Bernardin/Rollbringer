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
				h.renderErr(w, r, http.StatusUnauthorized, &domain.ProblemDetail{
					Type: domain.PDTypeUnauthorized,
				})
				return
			}

			h.renderErr(w, r, http.StatusInternalServerError, errors.Wrap(err, "cannot get SESSION_ID cookie"))
			return
		}
		sessionID, _ := uuid.Parse(stCookie.Value)

		session, err := h.Service.Authenticate(r.Context(), sessionID, true, r.Header.Get("CSRF-Token"))
		if err != nil {
			if domain.IsProblemDetail(err, domain.PDTypeUnauthorized) {
				h.err(w, r, &domain.ProblemDetail{
					Type: domain.PDTypeUnauthorized,
				})
				return
			}

			h.err(w, r, errors.Wrap(err, "cannot get session"))
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "session", session))
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) AuthenticatePage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get the session-ID from the cookie.
		stCookie, err := r.Cookie("SESSION_ID")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/users/login", http.StatusTemporaryRedirect)
				return
			}

			h.renderErr(w, r, http.StatusInternalServerError, errors.Wrap(err, "cannot get SESSION_ID cookie"))
			return
		}
		sessionID, _ := uuid.Parse(stCookie.Value)

		session, err := h.Service.Authenticate(r.Context(), sessionID, false, r.Header.Get("CSRF-Token"))
		if err != nil {
			if domain.IsProblemDetail(err, domain.PDTypeUnauthorized) {
				http.Redirect(w, r, "/users/login", http.StatusTemporaryRedirect)
				return
			}

			h.err(w, r, errors.Wrap(err, "cannot get session"))
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "session", session))
		next.ServeHTTP(w, r)
	})
}

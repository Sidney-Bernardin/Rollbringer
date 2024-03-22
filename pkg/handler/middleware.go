package handler

import (
	"context"
	"net/http"
	"rollbringer/pkg/domain"

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

func (h *Handler) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get the session-ID cookie.
		stCookie, err := r.Cookie("SESSION_ID")
		if err != nil {
			if err == http.ErrNoCookie {
				h.renderErr(w, r, http.StatusUnauthorized, &domain.ProblemDetail{
					Type: domain.PDTypeUnauthorized,
				})
				return
			}

			h.renderErr(w, r, http.StatusInternalServerError, errors.Wrap(err, "cannot get CSRF_Token cookie"))
			return
		}

		// Get the session.
		session, err := h.Service.GetSession(r.Context(), stCookie.Value)
		if err != nil {
			h.err(w, r, errors.Wrap(err, "cannot get session"))
			return
		}

		// Verify the CSRF-Token.
		if session.CSRFToken != r.Header.Get("CSRF-Token") {
			h.renderErr(w, r, http.StatusUnauthorized, &domain.ProblemDetail{
				Type: domain.PDTypeUnauthorized,
			})
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "session", session))
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) LightAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get the session-ID cookie.
		stCookie, err := r.Cookie("SESSION_ID")
		if err != nil {
			if err == http.ErrNoCookie {
				next.ServeHTTP(w, r)
				return
			}

			h.renderErr(w, r, http.StatusInternalServerError, errors.Wrap(err, "cannot get CSRF_Token cookie"))
			return
		}

		// Get the session.
		session, err := h.Service.GetSession(r.Context(), stCookie.Value)
		if err != nil && domain.IsProblemDetail(err, domain.PDTypeUnauthorized) {
			h.err(w, r, errors.Wrap(err, "cannot get session"))
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "session", session))
		next.ServeHTTP(w, r)
	})
}

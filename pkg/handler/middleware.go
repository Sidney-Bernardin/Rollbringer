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
				h.err(w, domain.ErrUnauthorized, http.StatusUnauthorized, 0)
				return
			}

			err = errors.Wrap(err, "cannot get CSRF_Token cookie")
			h.err(w, err, http.StatusInternalServerError, 0)
			return
		}

		// Get the session.
		session, err := h.Service.GetSession(r.Context(), stCookie.Value)
		if err != nil {
			h.domainErr(w, errors.Wrap(err, "cannot get session"))
			return
		}

		// Verify the CSRF-Token.
		if session.CSRFToken != r.Header.Get("CSRF-Token") {
			h.err(w, domain.ErrUnauthorized, http.StatusUnauthorized, 0)
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

			err = errors.Wrap(err, "cannot get CSRF_Token cookie")
			h.err(w, err, http.StatusInternalServerError, 0)
			return
		}

		// Get the session.
		session, err := h.Service.GetSession(r.Context(), stCookie.Value)
		if err != nil && errors.Cause(err) != domain.ErrUnauthorized {
			err = errors.Wrap(err, "cannot get session")
			h.err(w, err, http.StatusInternalServerError, 0)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "session", session))
		next.ServeHTTP(w, r)
	})
}

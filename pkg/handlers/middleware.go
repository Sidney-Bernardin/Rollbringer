package handlers

import (
	"net/http"
	"rollbringer/pkg/domain"
)

func (h *Handler) MWLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Logger.InfoContext(r.Context(), "New request",
			"method", r.Method,
			"url", r.URL.String())

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) MWAuthenticate(required, checkCSRF, redirect bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var (
				state = h.State(r)
				ctx   = r.Context()
			)

			session, err := h.authenticate(r, checkCSRF)
			if err != nil {
				h.Err(w, r, domain.Wrap(err, "cannot authenticate", nil))
				return
			}

			if session == nil && required {
				if redirect {
					http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				} else {
					h.Err(w, r, domain.UserErr(ctx, domain.UsrErrTypeUnauthorized, "You're unauthorized!", nil))
				}
				return
			}

			state["session"] = session
			next.ServeHTTP(w, r)
		})
	}
}

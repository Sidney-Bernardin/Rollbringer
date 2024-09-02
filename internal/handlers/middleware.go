package handlers

import (
	"context"
	"net/http"

	"rollbringer/internal"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (h *Handler) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Log the request.
		h.Logger.Info("New request",
			"method", r.Method,
			"url", r.URL.String(),
		)

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) Instance(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), internal.CtxKeyInstance, r.URL.Path))
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx = r.Context()

		cookie, err := r.Cookie("SESSION_ID")
		if err != nil {
			if err == http.ErrNoCookie {
				h.Err(w, r, internal.NewProblemDetail(ctx, internal.PDOpts{
					Type: internal.PDTypeUnauthorized,
				}))
				return
			}

			h.Err(w, r, errors.Wrap(err, "cannot get SESSION_ID cookie"))
			return
		}
		sessionID, _ := uuid.Parse(cookie.Value)

		session, err := h.svc.Authenticate(ctx, sessionID, r.Header.Get("CSRF-Token"))
		if err != nil {
			if internal.IsDetailed(err, internal.PDTypeUnauthorized) {
				h.Err(w, r, internal.NewProblemDetail(ctx, internal.PDOpts{
					Type: internal.PDTypeUnauthorized,
				}))
				return
			}

			h.Err(w, r, errors.Wrap(err, "cannot authenticate"))
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "session", session))
		next.ServeHTTP(w, r)
	})
}

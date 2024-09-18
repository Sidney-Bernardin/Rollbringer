package handlers

import (
	"context"
	"net/http"

	"rollbringer/internal"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (h *BaseHandler) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Log the request.
		h.Logger.Info("New request",
			"method", r.Method,
			"url", r.URL.String(),
		)

		next.ServeHTTP(w, r)
	})
}

func (h *BaseHandler) Instance(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), internal.CtxKeyInstance, r.URL.Path))
		next.ServeHTTP(w, r)
	})
}

func (h *BaseHandler) Authenticate(next http.Handler) http.Handler {
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

		session, err := h.BaseService.Authenticate(ctx, sessionID, r.Header.Get("CSRF-Token"))
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

		r = r.WithContext(context.WithValue(r.Context(), internal.CtxKeySession, session))
		next.ServeHTTP(w, r)
	})
}

func (h *BaseHandler) GetSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx = r.Context()

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

		session, err := h.BaseService.GetSession(ctx, sessionID, "session-all")
		if err != nil {
			if internal.IsDetailed(err, internal.PDTypeUnauthorized) {
				http.Redirect(w, r, "/users/login", http.StatusTemporaryRedirect)
				return
			}

			h.Err(w, r, errors.Wrap(err, "cannot get session"))
			return
		}

		r = r.WithContext(context.WithValue(ctx, internal.CtxKeySession, session))
		next.ServeHTTP(w, r)
	})
}

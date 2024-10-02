package handlers

import (
	"context"
	"net/http"

	"rollbringer/internal"

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

func (h *BaseHandler) Authenticate(sessionView internal.SessionView, checkCSRFToken bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()

			session, err := h.authenticate(r, sessionView, checkCSRFToken)
			if err != nil {
				if internal.IsDetailed(err, internal.PDTypeUnauthorized) {
					if checkCSRFToken {
						h.Err(w, r, internal.NewProblemDetail(ctx, internal.PDOpts{Type: internal.PDTypeUnauthorized}))
					} else {
						http.Redirect(w, r, "pages/login", http.StatusTemporaryRedirect)
					}
					return
				}

				h.Err(w, r, errors.Wrap(err, "cannot authenticate"))
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), internal.CtxKeySession, session))
			next.ServeHTTP(w, r)
		})
	}
}

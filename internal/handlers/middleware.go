package handlers

import (
	"context"
	"net/http"

	"rollbringer/internal"
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

package handlers

import (
	"net/http"
)

func (h *Handler) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Logger.InfoContext(r.Context(), "New request",
			"method", r.Method,
			"url", r.URL.String())

		next.ServeHTTP(w, r)
	})
}

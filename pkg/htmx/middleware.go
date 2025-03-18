package htmx

import (
	"net/http"
)

func (h *handler) mwInstance(instance string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.state(r)["instance"] = instance
		})
	}
}

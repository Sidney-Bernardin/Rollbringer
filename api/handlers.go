package api

import (
	"net/http"
)

func (a *api) handlePage(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := a.tmpl.ExecuteTemplate(w, name, nil); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

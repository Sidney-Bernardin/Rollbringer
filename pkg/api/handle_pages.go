package api

import (
	"fmt"
	"math/rand"
	"net/http"
)

func (a *api) handlePage(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]any{
			"name":     name,
			"oauthURL": a.oauthConfig.AuthCodeURL(fmt.Sprint(rand.Int())),
		}

		if err := a.tmpl.ExecuteTemplate(w, "page.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

func (a *api) handleGamePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]any{"name": "DND"}

		if err := a.tmpl.ExecuteTemplate(w, "page.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

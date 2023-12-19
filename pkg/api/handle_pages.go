package api

import (
	"net/http"
	"rollbringer/pkg/views"
)

func (api *API) HandleHomePage(w http.ResponseWriter, r *http.Request) {
	views.Home().Render(r.Context(), w)
}

func (api *API) HandleGamePage(w http.ResponseWriter, r *http.Request) {
	views.DND().Render(r.Context(), w)
}

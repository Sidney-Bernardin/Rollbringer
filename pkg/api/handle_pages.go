package api

import (
	"net/http"
	"rollbringer/pkg/models"
	"rollbringer/pkg/views"
)

func (api *API) HandleHomePage(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*models.Session)
	api.render(w, r, views.Home(session), http.StatusOK)
}

func (api *API) HandleGamePage(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*models.Session)
	api.render(w, r, views.DND(session), http.StatusOK)
}

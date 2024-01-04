package api

import (
	"net/http"
	"rollbringer/pkg/models"
	"rollbringer/pkg/views/pages/game"
	"rollbringer/pkg/views/pages/home"
)

func (api *API) HandleHomePage(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*models.Session)
	api.render(w, r, home.Home(session), http.StatusOK)
}

func (api *API) HandleGamePage(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*models.Session)
	api.render(w, r, game.DND(session), http.StatusOK)
}

package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/views/layouts"
)

func (api *API) HandlePlayPage(w http.ResponseWriter, r *http.Request) {

	// Get the game.
	game, err := api.Service.GetGame(r.Context(), chi.URLParam(r, "game_id"))
	if err != nil && err != domain.ErrGameNotFound {
		api.domainErr(w, errors.Wrap(err, "cannot get game"))
		return
	}
	giveToRequest(r, "game", game)

	// Check if the user is logged in by getting the session. If the user is
	// logged out, render the page early.
	session, _ := r.Context().Value("session").(*domain.Session)
	if session == nil {
		api.render(w, r, layouts.Play(), http.StatusOK)
		return
	}

	// Get the user.
	user, err := api.Service.GetUser(r.Context(), session.UserID.String())
	if err != nil {
		api.domainErr(w, errors.Wrap(err, "cannot get user"))
		return
	}
	giveToRequest(r, "user", user)

	// Get the user's games.
	games, err := api.Service.GetGamesFromUser(r.Context(), session.UserID.String())
	if err != nil {
		api.domainErr(w, errors.Wrap(err, "cannot get user's games"))
		return
	}
	giveToRequest(r, "games", games)

	api.render(w, r, layouts.Play(), http.StatusOK)
}

func (api *API) HandlePlayWS(conn *websocket.Conn) {

}

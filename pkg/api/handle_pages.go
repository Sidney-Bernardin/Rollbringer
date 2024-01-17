package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"rollbringer/pkg/models"
	database "rollbringer/pkg/repositories/database"
	"rollbringer/pkg/views/layouts"
)

func (api *API) HandleGamePage(w http.ResponseWriter, r *http.Request) {

	// Get the game-ID from the URL.
	gameID, err := uuid.Parse(r.URL.Query().Get("g"))
	if err == nil {

		// Get the game from the database.
		game, err := api.DB.GetGame(r.Context(), gameID)
		if err != nil && err != database.ErrGameNotFound {
			api.dbErr(w, errors.Wrap(err, "cannot get game"))
			return
		}
		giveToRequest(r, "game", game)
	}

	// Check if the user is logged in by getting the session. If the user is 
	// logged out, render the page early.
	session, _ := r.Context().Value("session").(*models.Session)
	if session == nil {
		api.renderHTTP(w, r, layouts.Game(), http.StatusOK)
		return
	}

	// Get the user from the database.
	user, err := api.DB.GetUser(r.Context(), session.UserID)
	if err != nil {
		api.dbErr(w, errors.Wrap(err, "cannot get user"))
		return
	}
	giveToRequest(r, "user", user)

	// Get user's games from the database.
	games, err := api.DB.GetGames(r.Context(), session.UserID)
	if err != nil {
		api.dbErr(w, errors.Wrap(err, "cannot get games"))
		return
	}
	giveToRequest(r, "games", games)

	api.renderHTTP(w, r, layouts.Game(), http.StatusOK)
}

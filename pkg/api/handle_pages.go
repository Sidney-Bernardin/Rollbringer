package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"rollbringer/pkg/models"
	"rollbringer/pkg/views/layouts"
)

func (api *API) HandleGamePage(w http.ResponseWriter, r *http.Request) {

	gameID, err := uuid.Parse(r.URL.Query().Get("g"))
	if err == nil {

		game, err := api.DB.GetGame(r.Context(), gameID)
		if err != nil {
			err = errors.Wrap(err, "cannot get game")
			api.err(w, r, err, http.StatusInternalServerError)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "game", game))
	}

	session, _ := r.Context().Value("session").(*models.Session)
	if session == nil {
		api.render(w, r, layouts.Game(), http.StatusOK)
		return
	}

	user, err := api.DB.GetUser(r.Context(), session.UserID)
	if err != nil {
		err = errors.Wrap(err, "cannot get user")
		api.err(w, r, err, http.StatusInternalServerError)
		return
	}
	r = r.WithContext(context.WithValue(r.Context(), "user", user))

	games, err := api.DB.GetGames(r.Context(), session.UserID)
	if err != nil {
		err = errors.Wrap(err, "cannot get games")
		api.err(w, r, err, http.StatusInternalServerError)
		return
	}
	r = r.WithContext(context.WithValue(r.Context(), "games", games))

	api.render(w, r, layouts.Game(), http.StatusOK)
}

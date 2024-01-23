package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/views/components"
)

func (api *API) HandleCreateGame(w http.ResponseWriter, r *http.Request) {

	session, _ := r.Context().Value("session").(*domain.Session)

	// Create a game.
	gameID, title, err := api.Service.CreateGame(r.Context(), session.UserID.String())
	if err != nil {
		api.domainErr(w, errors.Wrap(err, "cannot create game"))
		return
	}

	// Respond with a button for the new game.
	api.render(w, r, components.GameButton(gameID, title), http.StatusOK)
}

func (api *API) HandleDeleteGame(w http.ResponseWriter, r *http.Request) {

	session, _ := r.Context().Value("session").(*domain.Session)

	// Delete the game.
	if err := api.Service.DeleteGame(r.Context(), chi.URLParam(r, "game_id"), session.UserID.String()); err != nil {
		api.domainErr(w, errors.Wrap(err, "cannot delete game"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

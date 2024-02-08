package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/views/components/navigation"
)

func (api *API) handleCreateGame(w http.ResponseWriter, r *http.Request) {

	session, _ := r.Context().Value("session").(*domain.Session)

	// Create a game.
	gameID, title, err := api.service.CreateGame(r.Context(), session.UserID)
	if err != nil {
		api.domainErr(w, errors.Wrap(err, "cannot create game"))
		return
	}

	// Respond with a GameButton component.
	api.render(w, r, navigation.GameButton(gameID, title), http.StatusOK)
}

func (api *API) handleDeleteGame(w http.ResponseWriter, r *http.Request) {

	session, _ := r.Context().Value("session").(*domain.Session)

	// Delete the game.
	if err := api.service.DeleteGame(r.Context(), chi.URLParam(r, "game_id"), session.UserID); err != nil {
		api.domainErr(w, errors.Wrap(err, "cannot delete game"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

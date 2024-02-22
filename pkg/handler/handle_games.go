package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/views/components/navigation"
)

func (h *Handler) HandleCreateGame(w http.ResponseWriter, r *http.Request) {

	session, _ := r.Context().Value("session").(*domain.Session)

	// Create a game.
	game, err := h.Service.CreateGame(r.Context(), session)
	if err != nil {
		h.domainErr(w, errors.Wrap(err, "cannot create game"))
		return
	}

	// Respond with a GameButton component.
	h.render(w, r, navigation.GameButton(game), http.StatusOK)
}

func (h *Handler) HandleGetGames(w http.ResponseWriter, r *http.Request) {

	session, _ := r.Context().Value("session").(*domain.Session)

	// Get the games.
	games, err := h.Service.GetGames(r.Context(), session.UserID)
	if err != nil {
		h.domainErr(w, errors.Wrap(err, "cannot get games"))
		return
	}

	// Respond with a Games component.
	h.render(w, r, navigation.GameButtons(games), http.StatusOK)
}

func (h *Handler) HandleGetRolls(w http.ResponseWriter, r *http.Request) {

	// Get the rolls.
	rolls, err := h.Service.GetRolls(r.Context(), chi.URLParam(r, "game_id"))
	if err != nil {
		h.domainErr(w, errors.Wrap(err, "cannot get rolls"))
		return
	}

	// Respond with a Roll component.
	h.render(w, r, navigation.Rolls(rolls), http.StatusOK)
}
func (h *Handler) HandleDeleteGame(w http.ResponseWriter, r *http.Request) {

	session, _ := r.Context().Value("session").(*domain.Session)

	// Delete the game.
	if err := h.Service.DeleteGame(r.Context(), chi.URLParam(r, "game_id"), session.UserID); err != nil {
		h.domainErr(w, errors.Wrap(err, "cannot delete game"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

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
	gameID, title, err := h.Service.CreateGame(r.Context(), session.UserID)
	if err != nil {
		h.domainErr(w, errors.Wrap(err, "cannot create game"))
		return
	}

	// Respond with a GameButton component.
	h.render(w, r, navigation.GameButton(gameID, title), http.StatusOK)
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

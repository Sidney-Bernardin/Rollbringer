package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/views/components/navigation"
)

func (h *Handler) HandleCreateGame(w http.ResponseWriter, r *http.Request) {

	var session, _ = r.Context().Value("session").(*domain.Session)

	if session == nil {
		h.err(w, r, &domain.ProblemDetail{
			Type: domain.PDTypeUnauthorized,
		})
		return
	}

	game := &domain.Game{
		Name: "abc",
	}

	if err := h.Service.CreateGame(r.Context(), session, game); err != nil {
		h.err(w, r, errors.Wrap(err, "cannot create game"))
		return
	}

	h.render(w, r, http.StatusOK, navigation.HostedGameTableRow(game))
}

func (h *Handler) HandleDeleteGame(w http.ResponseWriter, r *http.Request) {

	var (
		session, _ = r.Context().Value("session").(*domain.Session)
		gameID, _  = uuid.Parse(chi.URLParam(r, "game_id"))
	)

	if session == nil {
		h.err(w, r, &domain.ProblemDetail{
			Type: domain.PDTypeUnauthorized,
		})
		return
	}

	if err := h.Service.DeleteGame(r.Context(), session, gameID); err != nil {
		h.err(w, r, errors.Wrap(err, "cannot delete game"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

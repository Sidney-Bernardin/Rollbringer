package games

import (
	"log/slog"
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"rollbringer/internal"
	"rollbringer/internal/config"
	"rollbringer/internal/handlers"
	service "rollbringer/internal/services/games"
	"rollbringer/internal/views/pages/play"
)

type handler struct {
	*handlers.BaseHandler

	svc service.Service
}

func NewHandler(cfg *config.Config, logger *slog.Logger, svc service.Service) http.Handler {
	h := &handler{
		BaseHandler: &handlers.BaseHandler{
			Config:      cfg,
			Logger:      logger,
			Router:      chi.NewRouter(),
			BaseService: svc,
		},
		svc: svc,
	}

	h.Router.Use(h.Log, h.Instance, h.Authenticate)
	h.Router.Post("/", h.HandleCreateGame)
	h.Router.Delete("/{game_id}", h.HandleDeleteGame)

	return h
}

func (h *handler) HandleCreateGame(w http.ResponseWriter, r *http.Request) {

	var (
		session, _ = r.Context().Value("session").(*internal.Session)
		game       = &internal.Game{
			Name: r.FormValue("name"),
		}
	)

	if err := h.svc.CreateGame(r.Context(), session, game); err != nil {
		h.Err(w, r, errors.Wrap(err, "cannot create game"))
		return
	}

	h.Render(w, r, http.StatusCreated, play.HostedGameRow(game))
}

func (h *handler) HandleDeleteGame(w http.ResponseWriter, r *http.Request) {

	var (
		session, _ = r.Context().Value("session").(*internal.Session)
		gameID, _  = uuid.Parse(chi.URLParam(r, "game_id"))
	)

	if err := h.svc.DeleteGame(r.Context(), session, gameID); err != nil {
		h.Err(w, r, errors.Wrap(err, "cannot delete game"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

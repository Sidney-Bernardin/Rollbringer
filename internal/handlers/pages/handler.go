package pages

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"rollbringer/internal"
	"rollbringer/internal/config"
	"rollbringer/internal/handlers"
	"rollbringer/internal/services/pages"
	html "rollbringer/internal/views/pages"
)

type handler struct {
	*handlers.BaseHandler

	svc pages.Service
}

func NewHandler(cfg *config.Config, logger *slog.Logger, svc pages.Service) *handler {
	h := &handler{
		BaseHandler: &handlers.BaseHandler{
			Config:      cfg,
			Logger:      logger,
			Router:      chi.NewRouter(),
			BaseService: svc,
		},
		svc: svc,
	}

	h.Router.Use(h.Log, h.Instance)
	h.Router.Get("/login", h.handleLogin)
	h.Router.With(h.Authenticate(internal.SessionViewPage, false)).Get("/", h.handlePlay)

	return h
}

func (h *handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	h.Render(w, r, http.StatusOK, html.Login())
}

func (h *handler) handlePlay(w http.ResponseWriter, r *http.Request) {
	var session, _ = r.Context().Value(internal.CtxKeySession).(*internal.Session)

	var game *internal.Game
	if gameID := internal.OptionalUUID(r.URL.Query().Get("g")); gameID != nil {
		game = &internal.Game{ID: *gameID}
	}

	page := &internal.PlayPage{
		Session: session,
		Game:    game,
	}

	if err := h.svc.PlayPage(r.Context(), page); err != nil {
		h.Err(w, r, errors.Wrap(err, "cannot get play-page"))
		return
	}

	h.Render(w, r, http.StatusOK, html.Play(page))
}

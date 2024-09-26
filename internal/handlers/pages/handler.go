package pages

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"rollbringer/internal/config"
	"rollbringer/internal/handlers"
	"rollbringer/internal/views/pages/play"
)

type handler struct {
	*handlers.BaseHandler
}

func NewHandler(cfg *config.Config, logger *slog.Logger) *handler {
	h := &handler{
		BaseHandler: &handlers.BaseHandler{
			Config:      cfg,
			Logger:      logger,
			Router:      chi.NewRouter(),
		},
	}

	h.Router.Use(h.Log, h.Instance)
	h.Router.Get("/login", h.handleLogin)
	h.Router.Get("/", h.handlePlay)

	return h
}

func (h *handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	h.Render(w, r, http.StatusOK, play.Login())
}

func (h *handler) handlePlay(w http.ResponseWriter, r *http.Request) {
	h.Render(w, r, http.StatusOK, play.Play())
}

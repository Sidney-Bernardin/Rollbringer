package pages

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"rollbringer/internal"
	"rollbringer/internal/config"
	"rollbringer/internal/handlers"
	service "rollbringer/internal/services/pages"
	"rollbringer/internal/views/pages/play"
)

type handler struct {
	*handlers.Handler

	svc service.Service
}

func NewHandler(cfg *config.Config, logger *slog.Logger, service service.Service) *handler {
	h := &handler{
		Handler: &handlers.Handler{
			Config: cfg,
			Logger: logger.With("component", "pages_handler"),
			Router: chi.NewRouter(),
		},
		svc: service,
	}

	h.Router.Use(h.Log, h.Instance, h.Authenticate)
	h.Router.Get("/play", h.HandlePlayPage)

	return h
}

func (h *handler) HandlePlayPage(w http.ResponseWriter, r *http.Request) {

	var (
		session, _ = r.Context().Value("session").(*internal.Session)
		gameID, _  = uuid.Parse(r.URL.Query().Get("g"))
	)

	page, err := h.svc.PlayPage(r.Context(), session, gameID)
	if err != nil {
		h.Err(w, r, errors.Wrap(err, "cannot get play page"))
		return
	}
	r = r.WithContext(context.WithValue(r.Context(), "play_page", page))

	h.Render(w, r, http.StatusOK, play.Play())
}

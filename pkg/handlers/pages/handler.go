package handler

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/handlers"
	"rollbringer/pkg/handlers/pages/views/pages"
)

type pagesHandler struct {
	*handlers.Handler

	pubSub domain.PubSubRepository
}

func New(config *domain.Config, logger *slog.Logger, svc *domain.Service) http.Handler {
	h := &pagesHandler{
		Handler: &handlers.Handler{
			Config:  config,
			Logger:  logger,
			Router:  chi.NewRouter(),
			Service: svc,
		},
	}

	h.Router.Use(h.MWLog)
	h.Router.Handle("/static/*", http.StripPrefix("/static/", http.FileServerFS(os.DirFS("cmd/static"))))
	h.Router.With(h.MWAuthenticate(false, false, true)).Get("/", h.handleHomePage)

	return h
}

func (h *pagesHandler) handleHomePage(w http.ResponseWriter, r *http.Request) {
	var session = h.State(r)["session"].(*domain.Session)

	h.Respond(w, r, http.StatusOK, pages.HomePage(&pages.HomePageState{
		Session: session,
	}))
}

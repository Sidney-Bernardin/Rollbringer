package handler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/handlers"
)

type PagesHandler struct {
	*handlers.Handler

	pubSub domain.PubSubRepository
}

func New(config *domain.Config, logger *slog.Logger, svc *domain.Service) http.Handler {
	h := &PagesHandler{
		Handler: &handlers.Handler{
			Config:  config,
			Logger:  logger,
			Router:  chi.NewRouter(),
			Service: svc,
		},
	}

	h.Router.Use(h.Log)
	h.Router.Post("/", h.HomePage)

	return h
}

type HomePage struct {
	Title string `json:"title"`
}

func (h *PagesHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	h.Respond(w, r, http.StatusOK, &HomePage{
		Title: "Home | " + r.RemoteAddr,
	})
}

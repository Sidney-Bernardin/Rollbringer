package handler

import (
	"log/slog"
	"net/http"

	"rollbringer/internal/config"
	"rollbringer/internal/handlers"
	"rollbringer/internal/services/users"

	"github.com/go-chi/chi/v5"
)

type usersHandler struct {
	*handlers.Handler

	service service.UsersService
}

func New(cfg *config.Config, logger *slog.Logger, service service.UsersService) *usersHandler {
	h := &usersHandler{
		Handler: &handlers.Handler{
			Router: chi.NewRouter(),
			Config: cfg,
			Logger: logger,
		},
		service: service,
	}

	h.Router.Use(h.Log, h.Instance)
	h.Router.Get("/login", h.HandleLogin)
	h.Router.Get("/consent-callback", h.HandleConsentCallback)

	return h
}

func (h *usersHandler) HandleLogin(w http.ResponseWriter, r *http.Request)           {}
func (h *usersHandler) HandleConsentCallback(w http.ResponseWriter, r *http.Request) {}
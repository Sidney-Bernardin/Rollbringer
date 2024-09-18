package users

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"rollbringer/internal"
	"rollbringer/internal/config"
	"rollbringer/internal/handlers"
	"rollbringer/internal/services/users"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

type usersHandler struct {
	*handlers.BaseHandler

	svc users.Service
}

func NewHandler(cfg *config.Config, logger *slog.Logger, service users.Service) *usersHandler {
	h := &usersHandler{
		BaseHandler: &handlers.BaseHandler{
			Config: cfg,
			Logger: logger,
			Router: chi.NewRouter(),
		},
		svc: service,
	}

	h.Router.Use(h.Log, h.Instance)
	h.Router.Get("/login", h.handleLogin)
	h.Router.Get("/consent-callback", h.handleConsentCallback)

	return h
}

func (h *usersHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	consentURL, state, codeVerifier := h.svc.StartLogin()

	http.SetCookie(w, &http.Cookie{
		Name:     "STATE_AND_VERIFIER",
		Value:    state + "," + codeVerifier,
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
	})

	http.Redirect(w, r, consentURL, http.StatusTemporaryRedirect)
}

func (h *usersHandler) handleConsentCallback(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	cookie, err := r.Cookie("STATE_AND_VERIFIER")
	if err != nil {
		h.Err(w, r, internal.NewProblemDetail(ctx, internal.PDOpts{
			Type: internal.PDTypeUnauthorized,
		}))
		return
	}

	// Get the state and code-verifier from the cookie.
	state_and_verifier := strings.Split(cookie.Value, ",")
	if len(state_and_verifier) != 2 {
		h.Err(w, r, internal.NewProblemDetail(ctx, internal.PDOpts{
			Type: internal.PDTypeUnauthorized,
		}))
		return
	}

	session, err := h.svc.FinishLogin(ctx,
		state_and_verifier[0],
		r.FormValue("state"),
		r.FormValue("code"),
		state_and_verifier[1],
	)

	if err != nil {
		h.Err(w, r, errors.Wrap(err, "cannot finish login"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "SESSION_ID",
		Value:    session.ID.String(),
		Path:     h.Config.CookiePath,
		Expires:  time.Now().Add(15 * time.Minute),
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/pages", http.StatusTemporaryRedirect)
}

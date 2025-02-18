package handler

import (
	"net/http"
	"time"

	"rollbringer/pkg/domain"

	"golang.org/x/oauth2"
)

func (h *accountsHandler) handleOAuth(w http.ResponseWriter, r *http.Request) {

	var (
		state = h.State(r)
		ctx   = r.Context()

		oauthConfig = state["oauth_config"].(*oauth2.Config)
	)

	// Create state string.
	oauthState, err := domain.NewRandomString(ctx)
	if err != nil {
		h.Err(w, r, domain.Wrap(err, "cannot create oauth state", nil))
		return
	}

	// Set a new cookie for the state string.
	http.SetCookie(w, &http.Cookie{
		Name:     "OAUTH_STATE",
		Value:    oauthState,
		Expires:  time.Now().Add(h.Config.OAuthStateCookieTimeout),
		HttpOnly: true,
	})

	// Redirect to OAuth consent page.
	http.Redirect(w, r, oauthConfig.AuthCodeURL(oauthState), http.StatusTemporaryRedirect)
}

func (h *accountsHandler) handleSignup(w http.ResponseWriter, r *http.Request) {
	var state = h.State(r)

	if err := h.accountsSvc.Signup(r.Context(), state["user"].(*domain.User)); err != nil {
		h.Err(w, r, domain.Wrap(err, "cannot signin", nil))
		return
	}
}

func (h *accountsHandler) handleSignin(w http.ResponseWriter, r *http.Request) {
	var state = h.State(r)

	if err := h.accountsSvc.Signin(r.Context(), state["user"].(*domain.User)); err != nil {
		h.Err(w, r, domain.Wrap(err, "cannot signin", nil))
		return
	}
}

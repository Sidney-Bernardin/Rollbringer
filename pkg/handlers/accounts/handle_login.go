package handler

import (
	"net/http"
	"time"

	"rollbringer/pkg/domain"

	"golang.org/x/oauth2"
)

func (h *AccountsHandler) handleOAuth(w http.ResponseWriter, r *http.Request) {
	var oauthConfig = r.Context().Value("oauth_config").(*oauth2.Config)

	// Create state string.
	state, err := domain.NewRandomString(r.Context())
	if err != nil {
		h.Err(w, r, domain.Wrap(err, "cannot create oauth state", nil))
		return
	}

	// Set a new cookie for the state string.
	http.SetCookie(w, &http.Cookie{
		Name:     "OAUTH_STATE",
		Value:    state,
		Expires:  time.Now().Add(h.Config.OAuthStateCookieTimeout),
		HttpOnly: true,
	})

	// Redirect to OAuth consent page.
	http.Redirect(w, r, oauthConfig.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func (h *AccountsHandler) handleSignup(w http.ResponseWriter, r *http.Request) {
	var user = r.Context().Value("user").(*domain.User)

	if err := h.accountsSvc.Signup(r.Context(), user); err != nil {
		h.Err(w, r, domain.Wrap(err, "cannot signin", nil))
		return
	}
}

func (h *AccountsHandler) handleSignin(w http.ResponseWriter, r *http.Request) {
	var user = r.Context().Value("user").(*domain.User)

	if err := h.accountsSvc.Signin(r.Context(), user); err != nil {
		h.Err(w, r, domain.Wrap(err, "cannot signin", nil))
		return
	}
}

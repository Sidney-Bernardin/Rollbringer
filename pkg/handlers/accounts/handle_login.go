package handler

import (
	"net/http"
	"time"

	"rollbringer/pkg/domain"

	"golang.org/x/oauth2"
)

func (h *accountsHandler) handleOAuth(w http.ResponseWriter, r *http.Request) {

	var (
		ctx         = r.Context()
		oauthConfig = ctx.Value("oauth_config").(*oauth2.Config)
	)

	// Create state string.
	state, err := domain.NewRandomString(ctx)
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

func (h *accountsHandler) handleSignup(w http.ResponseWriter, r *http.Request) {

	var (
		ctx  = r.Context()
		user = ctx.Value("user").(*domain.User)
	)

	if err := h.accountsSvc.Signup(ctx, user); err != nil {
		h.Err(w, r, domain.Wrap(err, "cannot signin", nil))
		return
	}
}

func (h *accountsHandler) handleSignin(w http.ResponseWriter, r *http.Request) {

	var (
		ctx  = r.Context()
		user = ctx.Value("user").(*domain.User)
	)

	if err := h.accountsSvc.Signin(ctx, user); err != nil {
		h.Err(w, r, domain.Wrap(err, "cannot signin", nil))
		return
	}
}

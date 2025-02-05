package handler

import (
	"context"
	"net/http"
	"time"

	"rollbringer/pkg/domain"

	"golang.org/x/oauth2"
)

func (h *AccountsHandler) handleOAuth(oauthConfig *oauth2.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var ctx = r.Context()

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
}

func (h *AccountsHandler) handleLoginCallbackGoogle(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	user, err := h.accountsSvc.LoginWithGoogle(ctx, ctx.Value("token").(*oauth2.Token))
	if err != nil {
		h.Err(w, r, domain.Wrap(err, "cannot login with google", nil))
		return
	}

	*r = *r.WithContext(context.WithValue(ctx, "user", user))
}

func (h *AccountsHandler) handleLoginCallbackSpotify(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	user, err := h.accountsSvc.LoginWithSpotify(ctx, h.oauthSpotifyConfig, ctx.Value("token").(*oauth2.Token))
	if err != nil {
		h.Err(w, r, domain.Wrap(err, "cannot login with spotify", nil))
		return
	}

	*r = *r.WithContext(context.WithValue(ctx, "user", user))
}

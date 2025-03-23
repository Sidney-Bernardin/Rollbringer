package api

import (
	"net/http"
	"time"

	"github.com/pkg/errors"

	"rollbringer/src"
)

func (svr *server) handleOAuthStart(w http.ResponseWriter, r *http.Request) {

	// Create state string.
	oauthState, err := src.NewRandomString()
	if err != nil {
		h.Err(w, r, errors.Wrap(err, "cannot create oauth state", nil))
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

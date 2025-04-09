package api

import (
	"net/http"
	"time"

	"github.com/pkg/errors"

	"rollbringer/src"
)

func (svr *server) handleOAuthConsent() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var consentURL, state string
		switch r.PathValue("provider") {
		case "google":
			consentURL, state = svr.google.ConsentURL()
		case "spotify":
			consentURL, state = svr.spotify.ConsentURL()
		default:
			svr.err(w, r, &src.ExternalError{Type: externalErrorTypeInvalidProvider})
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "OAUTH_STATE",
			Value:    state,
			HttpOnly: true,
			Expires:  time.Now().Add(svr.config.OAuthCookieTimeout),
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "OAUTH_LOGIN_TYPE",
			Value:    r.URL.Query().Get("login-type"),
			HttpOnly: true,
			Expires:  time.Now().Add(svr.config.OAuthCookieTimeout),
		})

		http.Redirect(w, r, consentURL, http.StatusTemporaryRedirect)
	})
}

func (svr *server) handleOAuthCallback() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx = r.Context()

		state, err := r.Cookie("OAUTH_STATE")
		if err != nil || state.Value != r.FormValue("state") {
			svr.err(w, r, &src.ExternalError{Type: src.ExternalErrorTypeUnauthorized})
			return
		}

		loginType, err := r.Cookie("OAUTH_LOGIN_TYPE")
		if err != nil {
			svr.err(w, r, &src.ExternalError{Type: src.ExternalErrorTypeUnauthorized})
			return
		}
		newAccount := loginType.Value == "signup"

		var sessionID *src.UUID
		switch r.PathValue("provider") {
		case "google":
			sessionID, err = svr.accounts.GoogleLogin(ctx, r.FormValue("code"), newAccount)
		case "spotify":
			sessionID, err = svr.accounts.SpotifyLogin(ctx, r.FormValue("code"), newAccount)
		default:
			err = &src.ExternalError{Type: externalErrorTypeInvalidProvider}
		}

		if err != nil {
			svr.err(w, r, errors.Wrap(err, "cannot login"))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "SESSION_ID",
			Value:    sessionID.String(),
			Path:     "/",
			SameSite: http.SameSiteStrictMode,
			HttpOnly: true,
			Expires:  time.Now().Add(svr.config.SessionCookieTimeout),
		})

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	})
}

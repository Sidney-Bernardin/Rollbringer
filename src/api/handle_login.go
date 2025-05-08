package api

import (
	"net/http"
	"rollbringer/src/domain"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	externalErrorTypeInvalidProvider domain.ExternalErrorType = "invalid-provider"
)

func (svr *server) handleOAuthConsent() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Generate the consent URL and it's state for the chosen provider.
		var consentURL, state string
		switch r.PathValue("provider") {
		case "google":
			consentURL, state = svr.google.ConsentURL()
		case "spotify":
			consentURL, state = svr.spotify.ConsentURL()
		default:
			svr.err(w, r, &domain.ExternalError{Type: externalErrorTypeInvalidProvider})
			return
		}

		// Save the state as a cookie.
		http.SetCookie(w, &http.Cookie{
			Name:     "OAUTH_STATE",
			Value:    state,
			HttpOnly: true,
			Expires:  time.Now().Add(svr.config.OAuthCookieTimeout),
		})

		// Save the login-type as a cookie.
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

		// Validate the state.
		state, err := r.Cookie("OAUTH_STATE")
		if err != nil || state.Value != r.FormValue("state") {
			svr.err(w, r, &domain.ExternalError{Type: domain.ExternalErrorTypeUnauthorized})
			return
		}

		// Get the login-type.
		loginType, err := r.Cookie("OAUTH_LOGIN_TYPE")
		if err != nil {
			svr.err(w, r, &domain.ExternalError{Type: domain.ExternalErrorTypeUnauthorized})
			return
		}

		// Login with the chosen provider.
		var sessionID uuid.UUID
		newAccount := loginType.Value == "signup"
		switch r.PathValue("provider") {
		case "google":
			sessionID, err = svr.accounts.GoogleLogin(ctx, r.FormValue("code"), newAccount)
		case "spotify":
			sessionID, err = svr.accounts.SpotifyLogin(ctx, r.FormValue("code"), newAccount)
		default:
			err = &domain.ExternalError{Type: externalErrorTypeInvalidProvider}
		}

		if err != nil {
			svr.err(w, r, errors.Wrap(err, "cannot login"))
			return
		}

		// Save the session-ID as a cookie.
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

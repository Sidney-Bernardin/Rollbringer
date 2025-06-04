package http

import (
	"net/http"
	"time"

	"github.com/Sidney-Bernardin/Rollbringer/server"

	"github.com/pkg/errors"
)

func (api *API) handleBasicLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		sessionID server.UUID
		err       error

		username = r.FormValue("username")
		password = r.FormValue("password")
	)

	switch r.URL.Query().Get("type") {
	case "signup":
		sessionID, err = api.Service.BasicSignup(ctx, username, password)
		err = errors.Wrap(err, "cannot signup")
	case "signin":
		sessionID, err = api.Service.BasicSignin(ctx, username, password)
		err = errors.Wrap(err, "cannot signin")
	default:
		api.err(w, r, server.NewUserError(server.UserErrorTypeUnauthorized, "", nil))
		return
	}

	if err != nil {
		api.err(w, r, errors.WithStack(err))
		return
	}

	http.SetCookie(w, api.NewSessionCookie(sessionID))
	w.Header().Set("HX-Redirect", "/")
}

func (api *API) handleGoogleLogin(w http.ResponseWriter, r *http.Request) {

	var (
		state      = server.CreateRandomString()
		consentURL = api.Service.Google.OAuthConfig.AuthCodeURL(state)
	)

	http.SetCookie(w, &http.Cookie{
		Name:     "OAUTH_STATE",
		Value:    state,
		Expires:  time.Now().Add(10 * time.Second),
		HttpOnly: true,
		Secure:   true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "LOGIN_TYPE",
		Value:    r.URL.Query().Get("type"),
		Expires:  time.Now().Add(10 * time.Second),
		HttpOnly: true,
		Secure:   true,
	})

	w.Header().Set("HX-Redirect", consentURL)
}

func (api *API) handleGoogleLoginCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cState, err := r.Cookie("OAUTH_STATE")
	if err != nil || cState.Value != r.FormValue("state") {
		api.err(w, r, server.NewUserError(server.UserErrorTypeUnauthorized, "", nil))
		return
	}

	cLoginType, err := r.Cookie("LOGIN_TYPE")
	if err != nil {
		api.err(w, r, server.NewUserError(server.UserErrorTypeUnauthorized, "", nil))
		return
	}

	googleUser, err := api.Service.Google.GetGoogleUser(ctx, r.FormValue("code"))
	if err != nil {
		api.err(w, r, errors.Wrap(err, "cannot get google-user"))
		return
	}

	var sessionID server.UUID
	switch cLoginType.Value {
	case "signup":
		sessionID, err = api.Service.GoogleSignup(ctx, googleUser)
		err = errors.Wrap(err, "cannot signup")
	case "signin":
		sessionID, err = api.Service.GoogleSignin(ctx, googleUser)
		err = errors.Wrap(err, "cannot signin")
	default:
		api.err(w, r, server.NewUserError(server.UserErrorTypeUnauthorized, "", nil))
		return
	}

	if err != nil {
		api.err(w, r, errors.WithStack(err))
		return
	}

	http.SetCookie(w, api.NewSessionCookie(sessionID))
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (api *API) handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "SESSION_ID",
		Expires: time.Unix(0, 0),
	})
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

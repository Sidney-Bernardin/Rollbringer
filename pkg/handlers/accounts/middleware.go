package handler

import (
	"net/http"
	"time"

	"golang.org/x/oauth2"

	"rollbringer/pkg/domain"
)

func (h *accountsHandler) mwOAuthConfig(oauthConfig *oauth2.Config, redirectURL string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var state = h.State(r)

			redirectOAuthConfig := *oauthConfig
			redirectOAuthConfig.RedirectURL = redirectURL

			state["oauth_config"] = &redirectOAuthConfig
			next.ServeHTTP(w, r)
		})
	}
}

func (h *accountsHandler) mwOAuthCallback(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			state = h.State(r)
			ctx   = r.Context()

			oauthConfig = state["oauth_config"].(*oauth2.Config)
		)

		stateCookie, err := r.Cookie("OAUTH_STATE")
		if err != nil {
			if err == http.ErrNoCookie {
				h.Err(w, r, domain.UserErr(ctx, domain.UsrErrTypeUnauthorized, "You're unauthorized!", nil))
				return
			}

			h.Err(w, r, domain.Wrap(err, "cannot get OAUTH_STATE cookie", nil))
			return
		}

		if r.FormValue("state") != stateCookie.Value {
			h.Err(w, r, domain.UserErr(ctx, domain.UsrErrTypeUnauthorized, "You're unauthorized!", nil))
			return
		}

		token, err := oauthConfig.Exchange(ctx, r.FormValue("code"))
		if err != nil {
			h.Err(w, r, domain.UserErr(ctx, domain.UsrErrTypeUnauthorized, "You're unauthorized!", nil))
			return
		}

		state["token"] = token
		next.ServeHTTP(w, r)

		if err := r.Context().Err(); err != nil {
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "SESSION_ID",
			Value:    state["user"].(*domain.User).Session.ID.String(),
			Path:     h.Config.SessionCookiePath,
			Expires:  time.Now().Add(h.Config.SessionCookieTimeout),
			SameSite: http.SameSiteStrictMode,
			HttpOnly: true,
		})

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	})
}

func (h *accountsHandler) mwCreateGoogleUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			state = h.State(r)
			ctx   = r.Context()

			token = state["token"].(*oauth2.Token)
		)

		user, err := h.accountsSvc.NewGoogleUser(ctx, token)
		if err != nil {
			h.Err(w, r, domain.Wrap(err, "cannot create google-user", nil))
			return
		}

		state["user"] = user
		next.ServeHTTP(w, r)
	})
}

func (h *accountsHandler) mwCreateSpotifyUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			state = h.State(r)
			ctx   = r.Context()

			oauthConfig = state["oauth_config"].(*oauth2.Config)
			token       = state["token"].(*oauth2.Token)
		)

		user, err := h.accountsSvc.NewSpotifyUser(ctx, oauthConfig, token)
		if err != nil {
			h.Err(w, r, domain.Wrap(err, "cannot create spotify-user", nil))
			return
		}

		state["user"] = user
		next.ServeHTTP(w, r)
	})
}

package handler

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/oauth2"

	"rollbringer/pkg/domain"
)

func (h *AccountsHandler) mwOAuthConfig(oauthConfig *oauth2.Config, redirectURL string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			redirectOAuthConfig := *oauthConfig
			redirectOAuthConfig.RedirectURL = redirectURL

			*r = *r.WithContext(context.WithValue(r.Context(), "oauth_config", &redirectOAuthConfig))
			next.ServeHTTP(w, r)
		})
	}
}

func (h *AccountsHandler) mwOAuthCallback(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var oauthConfig = r.Context().Value("oauth_config").(*oauth2.Config)

		stateCookie, err := r.Cookie("OAUTH_STATE")
		if err != nil {
			if err == http.ErrNoCookie {
				h.Err(w, r, domain.UserErr(r.Context(), domain.UsrErrTypeUnauthorized, "You're unauthorized!", nil))
				return
			}

			h.Err(w, r, domain.Wrap(err, "cannot get OAUTH_STATE cookie", nil))
			return
		}

		if r.FormValue("state") != stateCookie.Value {
			h.Err(w, r, domain.UserErr(r.Context(), domain.UsrErrTypeUnauthorized, "You're unauthorized!", nil))
			return
		}

		token, err := oauthConfig.Exchange(r.Context(), r.FormValue("code"))
		if err != nil {
			h.Err(w, r, domain.UserErr(r.Context(), domain.UsrErrTypeUnauthorized, "You're unauthorized!", nil))
			return
		}

		*r = *r.WithContext(context.WithValue(r.Context(), "token", token))
		next.ServeHTTP(w, r)

		if err := r.Context().Err(); err != nil {
			return
		}

		var user, _ = r.Context().Value("user").(*domain.User)

		http.SetCookie(w, &http.Cookie{
			Name:     "SESSION_ID",
			Value:    user.Session.ID.String(),
			Path:     h.Config.SessionCookiePath,
			Expires:  time.Now().Add(h.Config.SessionCookieTimeout),
			SameSite: http.SameSiteStrictMode,
			HttpOnly: true,
		})

		http.Redirect(w, r, "/pages", http.StatusTemporaryRedirect)
	})
}

func (h *AccountsHandler) mwCreateGoogleUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token = r.Context().Value("token").(*oauth2.Token)

		user, err := h.accountsSvc.NewGoogleUser(r.Context(), token)
		if err != nil {
			h.Err(w, r, domain.Wrap(err, "cannot create google-user", nil))
			return
		}

		*r = *r.WithContext(context.WithValue(r.Context(), "user", user))
		next.ServeHTTP(w, r)
	})
}

func (h *AccountsHandler) mwCreateSpotifyUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			oauthConfig = r.Context().Value("oauth_config").(*oauth2.Config)
			token       = r.Context().Value("token").(*oauth2.Token)
		)

		user, err := h.accountsSvc.NewSpotifyUser(r.Context(), oauthConfig, token)
		if err != nil {
			h.Err(w, r, domain.Wrap(err, "cannot create spotify-user", nil))
			return
		}

		*r = *r.WithContext(context.WithValue(r.Context(), "user", user))
		next.ServeHTTP(w, r)
	})
}

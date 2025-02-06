package handler

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/oauth2"

	"rollbringer/pkg/domain"
)

func (h *AccountsHandler) mwOAuthCallback(oauthConfig *oauth2.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()

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

			*r = *r.WithContext(context.WithValue(ctx, "token", token))
			next.ServeHTTP(w, r)

			user, _ := r.Context().Value("user").(*domain.User)
			if user == nil {
				return
			}

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
}

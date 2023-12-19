package api

import (
	"context"
	"net/http"
	database "rollbringer/pkg/repositories/database"

	"github.com/pkg/errors"
)

func (api *API) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		stCookie, err := r.Cookie("Session_Token")
		if err != nil {

			if err == http.ErrNoCookie {
				next.ServeHTTP(w, r)
				return
			}

			err = errors.Wrap(err, "cannot get CSRF_Token cookie")
			// server error
			return
		}

		session, err := api.DB.GetSession(r.Context(), stCookie.Value)
		if err != nil {

			if err == database.ErrSessionNotFound {
				next.ServeHTTP(w, r)
				return
			}

			err = errors.Wrap(err, "cannot get session from db")
			// server error
			return
		}

		if session.CSRFToken == r.Header.Get("CSRF-Token") {
			r = r.WithContext(context.WithValue(r.Context(), "session", session))
		}

		next.ServeHTTP(w, r)
	})
}

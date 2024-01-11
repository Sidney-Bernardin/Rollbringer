package api

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	database "rollbringer/pkg/repositories/database"
)

func (api *API) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		stCookie, err := r.Cookie("SESSION_ID")
		if err != nil {

			if err == http.ErrNoCookie {
				api.err(w, r, errUnauthorized, http.StatusUnauthorized)
				return
			}

			err = errors.Wrap(err, "cannot get CSRF_Token cookie")
			api.err(w, r, err, http.StatusInternalServerError)
			return
		}

		session, err := api.DB.GetSession(r.Context(), stCookie.Value)
		if err != nil {

			if err == database.ErrSessionNotFound {
				api.err(w, r, errUnauthorized, http.StatusUnauthorized)
				return
			}

			err = errors.Wrap(err, "cannot get session from db")
			api.err(w, r, err, http.StatusInternalServerError)
			return
		}

		if session.CSRFToken.String() == r.Header.Get("CSRF-Token") {
			r = r.WithContext(context.WithValue(r.Context(), "session", session))
		}

		next.ServeHTTP(w, r)
	})
}

func (api *API) LightAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		stCookie, err := r.Cookie("SESSION_ID")
		if err != nil {

			if err == http.ErrNoCookie {
				next.ServeHTTP(w, r)
				return
			}

			err = errors.Wrap(err, "cannot get CSRF_Token cookie")
			api.err(w, r, err, http.StatusInternalServerError)
			return
		}

		session, err := api.DB.GetSession(r.Context(), stCookie.Value)
		if err != nil {

			if err == database.ErrSessionNotFound {
				next.ServeHTTP(w, r)
				return
			}

			err = errors.Wrap(err, "cannot get session from db")
			api.err(w, r, err, http.StatusInternalServerError)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "session", session))
		next.ServeHTTP(w, r)
	})
}

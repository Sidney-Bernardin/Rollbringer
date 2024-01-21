package api

import (
	"net/http"

	"github.com/pkg/errors"

	database "rollbringer/pkg/repositories/database"
)

func (api *API) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.Logger.Info().Str("url", r.URL.String()).Msg("New request")
		next.ServeHTTP(w, r)
	})

}

func (api *API) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get the session-ID cookie.
		stCookie, err := r.Cookie("SESSION_ID")
		if err != nil {
			if err == http.ErrNoCookie {
				api.err(w, errUnauthorized, http.StatusUnauthorized, 0)
				return
			}

			err = errors.Wrap(err, "cannot get CSRF_Token cookie")
			api.err(w, err, http.StatusInternalServerError, 0)
			return
		}

		// Get the session.
		session, err := api.DB.GetSession(r.Context(), stCookie.Value)
		if err != nil {
			api.dbErr(w, errors.Wrap(err, "cannot get session from db"))
			return
		}

		// Verify the CSRF-Token.
		if session.CSRFToken.String() != r.Header.Get("CSRF-Token") {
			api.err(w, errUnauthorized, http.StatusUnauthorized, 0)
			return
		}

		giveToRequest(r, "session", session)
		next.ServeHTTP(w, r)
	})
}

func (api *API) LightAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get the session-ID cookie.
		stCookie, err := r.Cookie("SESSION_ID")
		if err != nil {
			if err == http.ErrNoCookie {
				next.ServeHTTP(w, r)
				return
			}

			err = errors.Wrap(err, "cannot get CSRF_Token cookie")
			api.err(w, err, http.StatusInternalServerError, 0)
			return
		}

		// Get the session.
		session, err := api.DB.GetSession(r.Context(), stCookie.Value)
		if err != nil && err != database.ErrUnauthorized {
			err = errors.Wrap(err, "cannot get session from db")
			api.err(w, err, http.StatusInternalServerError, 0)
			return
		}

		giveToRequest(r, "session", session)
		next.ServeHTTP(w, r)
	})
}

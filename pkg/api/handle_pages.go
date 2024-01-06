package api

import (
	"context"
	"net/http"
	"rollbringer/pkg/models"
	database "rollbringer/pkg/repositories/database"
	"rollbringer/pkg/views/pages/game"
	"rollbringer/pkg/views/pages/home"

	"github.com/pkg/errors"
)

func (api *API) HandleHomePage(w http.ResponseWriter, r *http.Request) {

	session, _ := r.Context().Value("session").(*models.Session)
	if session == nil {
		api.render(w, r, home.Home(), http.StatusOK)
		return
	}

	user, err := api.DB.GetUser(r.Context(), session.UserID)
	if err != nil {
		if cause := errors.Cause(err); cause == database.ErrUserNotFound {
			api.renderError(w, r, cause, http.StatusNotFound)
			return
		}

		err = errors.Wrap(err, "cannot get user")
		api.renderError(w, r, err, http.StatusInternalServerError)
		return
	}
	r = r.WithContext(context.WithValue(r.Context(), "user", user))

	api.render(w, r, home.Home(), http.StatusOK)
}

func (api *API) HandleGamePage(w http.ResponseWriter, r *http.Request) {
	api.render(w, r, game.DND(nil), http.StatusOK)
}

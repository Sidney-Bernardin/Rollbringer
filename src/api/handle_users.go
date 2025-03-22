package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"rollbringer/src/domain/accounts"
)

func (svr *server) handleUserCreate(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	// Initialize the user.
	var user any
	switch r.URL.Query().Get("u") {
	case "info":
		user = accounts.UserInfo{}
	}

	err := svr.accounts.UserCreate(ctx, user, &accounts.ArgsUserCreate{
		Username: r.FormValue("username"),
	})

	if err != nil {
		svr.err(w, r, errors.Wrap(err, "cannot create user"))
		return
	}

	svr.respond(w, r, http.StatusOK, user)
}

func (svr *server) handleUserGet(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	// Initialize the user.
	var user any
	switch r.URL.Query().Get("v") {
	case "info":
		user = accounts.UserInfo{}
	}

	if err := svr.accounts.UserGetByUsername(ctx, &user, chi.URLParam(r, "username")); err != nil {
		svr.err(w, r, errors.Wrap(err, "cannot get user by username"))
		return
	}

	svr.respond(w, r, http.StatusOK, user)
}

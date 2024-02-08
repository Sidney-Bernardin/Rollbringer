package api

import (
	"net/http"

	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/views/pages"
)

func (api *API) handlePlayPage(w http.ResponseWriter, r *http.Request) {

	// Get the game.
	game, err := api.service.GetGame(r.Context(), r.URL.Query().Get("g"))
	if err != nil && errors.Cause(err) != domain.ErrGameNotFound {
		api.domainErr(w, errors.Wrap(err, "cannot get game"))
		return
	}
	giveToRequest(r, "game", game)

	// Check if the user is logged in by getting the session. If the user is
	// logged out, render the page early.
	session, _ := r.Context().Value("session").(*domain.Session)
	if session == nil {
		api.render(w, r, pages.Play(), http.StatusOK)
		return
	}

	// Get the user.
	user, err := api.service.GetUser(r.Context(), session.UserID)
	if err != nil {
		api.domainErr(w, errors.Wrap(err, "cannot get user"))
		return
	}
	giveToRequest(r, "user", user)

	// Get the user's games.
	games, err := api.service.GetGamesFromUser(r.Context(), session.UserID)
	if err != nil {
		api.domainErr(w, errors.Wrap(err, "cannot get user's games"))
		return
	}
	giveToRequest(r, "games", games)

	// Get the user's PDFs.
	pdfs, err := api.service.GetPDFs(r.Context(), session.UserID)
	if err != nil {
		api.domainErr(w, errors.Wrap(err, "cannot get pdfs"))
		return
	}
	giveToRequest(r, "pdfs", pdfs)

	api.render(w, r, pages.Play(), http.StatusOK)
}

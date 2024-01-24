package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/views/layouts"
)

func (api *API) HandlePlayPage(w http.ResponseWriter, r *http.Request) {

	// Get the game.
	game, err := api.service.GetGame(r.Context(), chi.URLParam(r, "game_id"))
	if err != nil && err != domain.ErrGameNotFound {
		api.domainErr(w, errors.Wrap(err, "cannot get game"))
		return
	}
	giveToRequest(r, "game", game)

	// Check if the user is logged in by getting the session. If the user is
	// logged out, render the page early.
	session, _ := r.Context().Value("session").(*domain.Session)
	if session == nil {
		api.render(w, r, layouts.Play(), http.StatusOK)
		return
	}

	// Get the user.
	user, err := api.service.GetUser(r.Context(), session.UserID.String())
	if err != nil {
		api.domainErr(w, errors.Wrap(err, "cannot get user"))
		return
	}
	giveToRequest(r, "user", user)

	// Get the user's games.
	games, err := api.service.GetGamesFromUser(r.Context(), session.UserID.String())
	if err != nil {
		api.domainErr(w, errors.Wrap(err, "cannot get user's games"))
		return
	}
	giveToRequest(r, "games", games)

	api.render(w, r, layouts.Play(), http.StatusOK)
}

func (api *API) HandlePlayWS(conn *websocket.Conn) {

	defer func() {
		if err := conn.Close(); err != nil {
			api.logger.Error().Stack().Err(err).Msg("Cannot close connection")
		}
	}()

	var (
		r = conn.Request()

		incomingChan = make(chan domain.GameEvent)
		outgoingChan = make(chan domain.GameEvent)
	)

	go api.service.Play(r.Context(), chi.URLParam(r, "game_id"), incomingChan, outgoingChan)

	go func() {
		for {

			var msg domain.GameEvent
			if err := websocket.JSON.Receive(conn, &msg); err != nil {

				if err == io.EOF {
					return
				}

				switch err.(type) {
				case *json.SyntaxError, *json.UnmarshalTypeError, *json.InvalidUnmarshalError:
					api.err(conn, err, 0, wsStatusUnsupportedData)
					return
				}

				api.err(conn, err, 0, wsStatusInternalError)
				return
			}

			incomingChan <- msg
		}
	}()

	for {
		select {

		case <-r.Context().Done():
			return

		case event := <-outgoingChan:
			if err := websocket.JSON.Send(conn, event); err != nil {
				api.err(conn, err, 0, wsStatusInternalError)
				return
			}
		}
	}
}

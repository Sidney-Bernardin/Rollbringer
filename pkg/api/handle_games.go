package api

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"

	"rollbringer/pkg/models"
	"rollbringer/pkg/repositories/database"
	"rollbringer/pkg/views/components"
)

func (a *API) HandleCreateGame(w http.ResponseWriter, r *http.Request) {

	session, _ := r.Context().Value("session").(*models.Session)

	gameID, title, err := a.DB.CreateGame(r.Context(), session.UserID)
	if err != nil {

		if err == database.ErrMaxGames {
			a.err(w, r, err, http.StatusForbidden)
			return
		}

		err = errors.Wrap(err, "cannot create game")
		a.err(w, r, err, http.StatusInternalServerError)
		return
	}

	a.render(w, r, components.GameButton(gameID, title), http.StatusOK)
}

func (a *API) HandleDeleteGame(w http.ResponseWriter, r *http.Request) {

	session, _ := r.Context().Value("session").(*models.Session)
	gameID, err := uuid.Parse(chi.URLParam(r, "game_id"))
	if err != nil {
		a.err(w, r, database.ErrGameNotFound, http.StatusNotFound)
		return
	}

	if err := a.DB.DeleteGame(r.Context(), session.UserID, gameID); err != nil {

		switch err {
		case database.ErrUnauthorized:
			a.err(w, r, err, http.StatusUnauthorized)
			return
		case database.ErrGameNotFound:
			a.err(w, r, err, http.StatusNotFound)
			return
		}

		err = errors.Wrap(err, "cannot delete game")
		a.err(w, r, err, http.StatusInternalServerError)
		return
	}
}

func (a *API) HandleJoinGame(conn *websocket.Conn) {
	req := conn.Request()

	res := components.HTMxAddTabs(
		components.TabPanelSelectorPlayMaterial,
		map[string]templ.Component{
			"Hoid": components.DNDCharacterSheet(),
			"Lee":  components.DNDCharacterSheet(),
		},
	)

	if err := res.Render(req.Context(), conn); err != nil {
		fmt.Println(err)
		return
	}

	for {
	}
}

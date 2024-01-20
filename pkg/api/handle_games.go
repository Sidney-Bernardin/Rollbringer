package api

import (
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

func (api *API) HandleCreateGame(w http.ResponseWriter, r *http.Request) {

	session, _ := r.Context().Value("session").(*models.Session)

	// Insert a new game.
	gameID, title, err := api.DB.InsertGame(r.Context(), session.UserID)
	if err != nil {
		api.dbErr(w, errors.Wrap(err, "cannot create game"))
		return
	}

	// Respond with a button for the new game
	api.renderHTTP(w, r, components.GameButton(gameID, title), http.StatusOK)
}

func (api *API) HandleDeleteGame(w http.ResponseWriter, r *http.Request) {

	session, _ := r.Context().Value("session").(*models.Session)

	// Get and parse the game-ID from the URL.
	gameID, err := uuid.Parse(chi.URLParam(r, "game_id"))
	if err != nil {
		api.dbErr(w, database.ErrGameNotFound)
		return
	}

	// Delete the game.
	if err := api.DB.DeleteGame(r.Context(), gameID, session.UserID); err != nil {
		api.dbErr(w, errors.Wrap(err, "cannot delete game"))
	}
}

func (api *API) HandleJoinGame(conn *websocket.Conn) {
	req := conn.Request()

	res := components.HTMXAddTabs(
		components.TabPanelSelectorPlayMaterial,
		map[string]templ.Component{
			"Hoid": components.DNDCharacterSheet(),
			"Lee":  components.DNDCharacterSheet(),
		},
	)

	api.renderWS(req.Context(), conn, res)

	for {
	}
}

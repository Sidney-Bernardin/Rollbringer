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

	// Get the session from the request's context.
	session, _ := r.Context().Value("session").(*models.Session)

	// Create the game.
	gameID, title, err := api.DB.CreateGame(r.Context(), session.UserID)
	if err != nil {
		api.dbErr(w, errors.Wrap(err, "cannot create game"))
		return
	}

	// Respond with a button for the new game
	api.renderHTTP(w, r, components.GameButton(gameID, title), http.StatusOK)
}

func (api *API) HandleDeleteGame(w http.ResponseWriter, r *http.Request) {

	// Get the session and game-ID.
	session, _ := r.Context().Value("session").(*models.Session)
	gameID, err := uuid.Parse(chi.URLParam(r, "game_id"))
	if err != nil {
		api.dbErr(w, database.ErrGameNotFound)
		return
	}

	// Delete the game.
	if err := api.DB.DeleteGame(r.Context(), session.UserID, gameID); err != nil {
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

package api

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"

	"rollbringer/pkg/models"
	"rollbringer/pkg/repositories/database"
	"rollbringer/pkg/views/components"
)

func (a *API) HandleCreateGame(w http.ResponseWriter, r *http.Request) {

	session, _ := r.Context().Value("session").(*models.Session)

	title, err := a.DB.CreateGame(r.Context(), session)
	if err != nil {
		if err == database.ErrMaxGames {
			a.renderError(w, r, err, http.StatusForbidden)
			return
		}

		err = errors.Wrap(err, "cannot create game")
		a.renderError(w, r, err, http.StatusInternalServerError)
		return
	}

	a.render(w, r, components.GameButton(title), http.StatusOK)
}

func (a *API) HandleJoinGame(conn *websocket.Conn) {
	req := conn.Request()
	chiCtx := chi.RouteContext(req.Context())

	_ = chiCtx.URLParam("id")

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

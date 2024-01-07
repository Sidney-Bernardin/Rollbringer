package api

import (
	"fmt"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"golang.org/x/net/websocket"

	"rollbringer/pkg/views/components"
)

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

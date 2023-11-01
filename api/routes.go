package api

import (
	"io/fs"
	"net/http"
)

func (a *api) doRoutes(staticFS fs.FS) {

	// Serve static files.
	a.router.Handle(
		"/static/*",
		http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))),
	)

	a.router.Get("/", a.handlePage("home.html"))
	a.router.Get("/game", a.handlePage("game.html"))
}

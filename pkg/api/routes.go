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

	a.router.Get("/", a.handlePage("Home"))
	a.router.Get("/game", a.handleGamePage())
	a.router.Get("/signin-callback", a.handleSigninCallback())
	a.router.Handle("/ws", a.handleWS())
}

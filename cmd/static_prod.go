//go:build prod
// +build prod

package main

import (
	"embed"
	"net/http"
)

//go:embed static
var staticFS embed.FS

func handleStaticDir() http.Handler {
	return http.FileServerFS(staticFS)
}

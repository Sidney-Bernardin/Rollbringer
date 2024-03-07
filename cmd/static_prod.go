//go:build prod

package main

import (
	"embed"
	"fmt"
	"net/http"
)

//go:embed static
var staticFS embed.FS

func handleStaticDir() http.Handler {
	fmt.Println("prod")
	return http.FileServerFS(staticFS)
}

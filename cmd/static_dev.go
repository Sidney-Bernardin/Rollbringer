//go:build !prod
// +build !prod

package main

import (
	"net/http"
	"os"
)

func handleStaticDir() http.Handler {
	return http.StripPrefix("/static/", http.FileServerFS(os.DirFS("cmd/static")))
}

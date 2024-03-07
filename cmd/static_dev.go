//go:build !prod
// +build !prod

package main

import (
	"fmt"
	"net/http"
	"os"
)

func handleStaticDir() http.Handler {
	fmt.Println("dev")
	return http.StripPrefix("/static/", http.FileServerFS(os.DirFS("cmd/static")))
}

package main

import (
	"rollbringer/server/domain/play"
	"rollbringer/server/htmx"
)

func main() {

	playSvc := play.NewService(nil, nil)
	handler := htmx.NewHandler(logger)
}

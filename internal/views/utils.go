package views

import (
	"context"
	"fmt"

	"rollbringer/internal"
	pages_service "rollbringer/internal/services/pages"
)

var S = fmt.Sprint
var F = fmt.Sprintf

func GetPlayPage(ctx context.Context) (page *pages_service.PlayPage) {
	if page, _ = ctx.Value("play_page").(*pages_service.PlayPage); page == nil {
		panic("play page is nil")
	}
	return page
}

func GetSession(ctx context.Context) (session *internal.Session) {
	if session, _ = ctx.Value("session").(*internal.Session); session == nil {
		panic("session is nil")
	}
	return session
}

package views

import (
	"fmt"
    "context"

	"rollbringer/pkg/domain"
)

var F = fmt.Sprintf

func GetPlayPage(ctx context.Context) (page *domain.PlayPage) {
	if page, _ = ctx.Value("play_page").(*domain.PlayPage); page == nil {
		panic("bad play page")
	}
	return page
}

func GetSession(ctx context.Context) (session *domain.Session) {
	if session, _ = ctx.Value("session").(*domain.Session); session == nil {
		panic("bad session")
	}
	return session
}

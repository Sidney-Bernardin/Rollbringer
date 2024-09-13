package pages

import "rollbringer/internal"

const CtxKey internal.CtxKey = "play_page"

type PlayPage struct {
	IsHost bool

	User *internal.User
	Game *internal.Game
}

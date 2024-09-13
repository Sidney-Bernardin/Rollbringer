package games

import (
	"rollbringer/internal"
	"rollbringer/internal/repositories/database"
)

type gamesSchema struct {
	*database.Database
}

func New(db *database.Database) internal.GamesSchema {
	return &gamesSchema{
		Database: db,
	}
}

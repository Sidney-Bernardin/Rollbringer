package users

import (
	"rollbringer/internal"
	"rollbringer/internal/repositories/database"
)

type usersSchema struct {
	*database.Database
}

func New(db *database.Database) (internal.UsersSchema) {
	return &usersSchema{
		Database: db,
	}
}

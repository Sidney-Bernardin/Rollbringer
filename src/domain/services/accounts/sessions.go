package accounts

import (
	"rollbringer/src"

	"github.com/google/uuid"
)

type Session struct {
	ID uuid.UUID

	UserID uuid.UUID
	User   *User

	CSRFToken CSRFToken
}

type CSRFToken string

func NewCSRFToken() CSRFToken {
	return CSRFToken(src.CreateRandomString())
}

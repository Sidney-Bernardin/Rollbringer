package models

import "rollbringer/src"

type Session struct {
	ID src.UUID `json:"id"`

	UserID src.UUID `json:"user_id"`
	User   *User    `json:"user"`

	CSRFToken CSRFToken `json:"csrf_token"`
}

type CSRFToken string

func NewCSRFToken() CSRFToken {
	return CSRFToken(src.CreateRandomString())
}

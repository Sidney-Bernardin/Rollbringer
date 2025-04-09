package models

import "rollbringer/src"

type Session struct {
	SessionID src.UUID `json:"session_id"`

	UserID src.UUID `json:"user_id"`
	User   *User    `json:"user"`

	CSRFToken CSRFToken `json:"csrf_token"`
}

type CSRFToken string

func NewCSRFToken() CSRFToken {
	return CSRFToken(src.CreateRandomString())
}

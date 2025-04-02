package accounts

import "rollbringer/src"

type CSRFToken string

func NewCSRFToken() CSRFToken {
	return CSRFToken(src.CreateRandomString())
}

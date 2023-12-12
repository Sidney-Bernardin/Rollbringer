package database

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	GoogleID string
	Username string
}

type Session struct {
	ID        uuid.UUID
	CSRFToken string
}

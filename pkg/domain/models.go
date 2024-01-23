package domain

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `db:"id,omitempty"`
	GoogleID string    `db:"google_id,omitempty"`
	Username string    `db:"username,omitempty"`
}

type Session struct {
	ID        uuid.UUID `db:"id,omitempty"`
	CSRFToken uuid.UUID `db:"csrf_token,omitempty"`
	UserID    uuid.UUID `db:"user_id,omitempty"`
}

type Game struct {
	ID     uuid.UUID
	HostID uuid.UUID
	Title  string
}

type GameEvent map[string]any

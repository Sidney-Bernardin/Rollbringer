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
	ID     uuid.UUID `db:"id,omitempty"`
	HostID uuid.UUID `db:"host_id,omitempty"`
	Title  string    `db:"title,omitempty"`
}

type PDF struct {
	ID      uuid.UUID `db:"id,omitempty"`
	OwnerID uuid.UUID `db:"owner_id,omitempty"`
	Name    string    `db:"name,omitempty"`
	Schema  string    `db:"schema,omitempty"`
	Content []byte    `db:"content,omitempty"`
}

type GameEvent map[string]any

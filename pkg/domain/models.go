package domain

import (
	"time"

	"github.com/google/uuid"
)

type Operation string

const (
	OperationError Operation = "ERROR"
)

type Event struct {
	Operation Operation `json:"operation"`
	Payload   any       `json:"payload"`
}

/////

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	Username string `json:"username" db:"username"`

	Session     *Session     `json:"session" db:"session"`
	GoogleUser  *GoogleUser  `json:"google_user" db:"google_user"`
	SpotifyUser *SpotifyUser `json:"spotify_user" db:"spotify_user"`
}

type GoogleUser struct {
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	GoogleID       string `json:"google_id" db:"google_id"`
	GivenName      string `json:"given_name" db:"given_name"`
	ProfilePicture string `json:"profile_picture" db:"profile_picture"`
}

type SpotifyUser struct {
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	SpotifyID      string `json:"spotify_id" db:"spotify_id"`
	DisplayName    string `json:"display_name" db:"display_name"`
	ProfilePicture string `json:"profile_picture" db:"profile_picture"`
}

type Session struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	CSRFToken string    `json:"csrf_token" db:"csrf_token"`
}

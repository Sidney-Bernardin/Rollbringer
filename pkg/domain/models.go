package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	GoogleID  *string `json:"google_id" db:"google_id"`
	SpotifyID *string `json:"spotify_id" db:"spotify_id"`

	Username       string `json:"username" db:"username"`
	ProfilePicture string `json:"profile_picture" db:"profile_picture"`

	GoogleUser  *GoogleUser  `json:"google_user,omitempty" db:"google_user"`
	SpotifyUser *SpotifyUser `json:"spotify_user,omitempty" db:"spotify_user"`
	Session     *Session     `json:"session,omitempty" db:"session"`
}

type GoogleUser struct {
	GoogleID  string    `json:"google_id" db:"google_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	GivenName      string `json:"given_name" db:"given_name"`
	Email          string `json:"email" db:"email"`
	ProfilePicture string `json:"profile_picture" db:"profile_picture"`
}

type SpotifyUser struct {
	SpotifyID string    `json:"spotify_id" db:"spotify_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	DisplayName    string  `json:"display_name" db:"display_name"`
	Email          string  `json:"email" db:"email"`
	ProfilePicture *string `json:"profile_picture" db:"profile_picture"`
}

type Session struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	CSRFToken string    `json:"csrf_token" db:"csrf_token"`

	User        *User        `json:"user,omitempty" db:"user"`
	GoogleUser  *GoogleUser  `json:"google_user,omitempty" db:"google_user"`
	SpotifyUser *SpotifyUser `json:"spotify_user,omitempty" db:"spotify_user"`
}

package domain

import (
	"time"
)

const (
	DomainErrorTypeUsernameInvalid UserErrorType = "username-invalid"
)

type User struct {
	ID UUID `json:"id" db:"id"`

	CreatedAt time.Time `json:"created_at,omitzero" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitzero" db:"updated_at"`

	GoogleID   *string     `json:"google_id,omitempty" db:"google_id"`
	GoogleUser *GoogleUser `json:"google_user,omitempty"`

	SpotifyID   *string      `json:"spotify_id,omitempty" db:"spotify_id"`
	SpotifyUser *SpotifyUser `json:"spotify_user,omitempty"`

	Rooms []*Room `json:"rooms,omitempty"`

	Username       Username `json:"username,omitempty" db:"username"`
	ProfilePicture string   `json:"profile_picture,omitempty" db:"profile_picture"`
}

type Username string

func ParseUsername(username string) (Username, error) {
	if len(username) < 5 || 32 < len(username) {
		return "", &UserError{
			Type:    DomainErrorTypeUsernameInvalid,
			Message: "Username must be between 5 and 32 characters.",
		}
	}

	return Username(username), nil
}

type GoogleUser struct {
	GoogleID string `json:"google_id" db:"google_id"`

	CreatedAt time.Time `json:"created_at,omitzero" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitzero" db:"updated_at"`

	GivenName      string `json:"given_name,omitempty" db:"given_name"`
	Email          string `json:"email,omitempty" db:"email"`
	ProfilePicture string `json:"profile_picture,omitempty"`
}

type SpotifyUser struct {
	SpotifyID string `json:"spotify_id" db:"spotify_id"`

	CreatedAt time.Time `json:"created_at,omitzero" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitzero" db:"updated_at"`

	DisplayName    string  `json:"display_name,omitempty" db:"display_name"`
	Email          string  `json:"email,omitempty" db:"email"`
	ProfilePicture *string `json:"profile_picture,omitempty"`
}

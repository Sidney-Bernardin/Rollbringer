package accounts

import "github.com/google/uuid"

type user struct {
	ID             uuid.UUID `db:"id"`
	GoogleID       *string   `db:"google_id"`
	SpotifyID      *string   `db:"spotify_id"`
	Username       string    `db:"username"`
	ProfilePicture string    `db:"profile_picture"`
}

package accounts

import (
	"github.com/google/uuid"
)

type googleUser struct {
	GoogleID uuid.UUID `db:"google_id"`

	GivenName      string `db:"given_name"`
	Email          string `db:"email"`
	ProfilePicture string `db:"profile_picture"`
}

const (
	qGoogleUserInsert = `
		INSERT INTO accounts.google_users (google_id, given_name, email, profile_picture)
		VALUES ($1, $2, $3, $4)`

	qGoogleUserUpdateByID = `
		UPDATE accounts.google_users 
			SET given_name = $2, email = $3, profile_picture = $4
			WHERE google_id = $1`
)

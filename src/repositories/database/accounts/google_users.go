package accounts

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	"rollbringer/src"
	"rollbringer/src/repositories/database"
	"rollbringer/src/services/accounts/models"
)

type googleUser struct {
	GoogleID string `db:"google_id"`

	GivenName      string `db:"given_name"`
	Email          string `db:"email"`
	ProfilePicture string `db:"profile_picture"`
}

func (db *accountsDatabase) GoogleSignup(ctx context.Context, user *models.User) (*src.UUID, error) {
	var sessionID = src.NewUUID()

	err := db.Transaction(ctx, func(tx pgx.Tx) error {

		// Insert the user and it's google-user.
		err := database.Insert(ctx, tx, `
			WITH inserted_google_user AS (
				INSERT INTO accounts.google_users (google_id, given_name, email, profile_picture)
				VALUES ($4, $5, $6, $7)
			)
			INSERT INTO accounts.users (id, google_id, username, profile_picture)
			VALUES ($1, $4, $2, $3)
		`,
			user.ID, user.Username, user.ProfilePicture,
			user.GoogleUser.GoogleID, user.GoogleUser.GivenName, user.GoogleUser.Email, user.GoogleUser.ProfilePicture)
		if err != nil {
			return errors.Wrap(err, "cannot insert user and google-user")
		}

		// Upsert a new session for the user.
		err = database.Insert(ctx, tx, `
			INSERT INTO accounts.sessions (id, user_id, csrf_token)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id) DO UPDATE SET
				id = EXCLUDED.id,
				csrf_token = EXCLUDED.csrf_token
		`,
			sessionID, user.ID, models.NewCSRFToken())
		return errors.Wrap(err, "cannot upsert session")
	})

	return &sessionID, errors.Wrap(err, "transaction failed")
}

func (db *accountsDatabase) GoogleSignin(ctx context.Context, googleUser *models.GoogleUser) (*src.UUID, error) {
	var sessionID = src.NewUUID()

	// Update the google-user.
	err := database.Update(ctx, db.Tx, `
		UPDATE accounts.google_users 
			SET given_name = $2, email = $3, profile_picture = $4
			WHERE google_id = $1
	`,
		googleUser.GoogleID, googleUser.GivenName, googleUser.Email, googleUser.ProfilePicture)
	if err != nil {
		return nil, errors.Wrap(err, "cannot update google-user by ID")
	}

	// Get google-user's user.
	user, err := database.Get[userRow](ctx, db.Tx, `
		SELECT users.id AS "users.id"
		FROM accounts.users
		WHERE users.google_id = $1
	`, googleUser.GoogleID)
	if err != nil {
		return nil, errors.Wrap(err, "cannot select user by google-ID")
	}

	// Upsert a new session for the user.
	err = database.Insert(ctx, db.Tx, `
		INSERT INTO accounts.sessions (id, user_id, csrf_token)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE SET
			id = EXCLUDED.id,
			csrf_token = EXCLUDED.csrf_token
	`,
		sessionID, user.ID, models.NewCSRFToken())
	return &sessionID, errors.Wrap(err, "cannot upsert session")
}

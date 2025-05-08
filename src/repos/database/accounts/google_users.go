package accounts

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	"rollbringer/src/domain/services/accounts"
	"rollbringer/src/repos/database"
)

type googleUser struct {
	GoogleID string `db:"google_id"`

	GivenName      string `db:"given_name"`
	Email          string `db:"email"`
	ProfilePicture string `db:"profile_picture"`
}

func (db *accountsDatabase) GoogleSignup(ctx context.Context, googleUser *accounts.GoogleUser, user *accounts.User) (sessionID uuid.UUID, err error) {
	err = db.Transaction(ctx, func(tx pgx.Tx) error {

		// Insert the user and it's google-user.
		err := database.Insert(ctx, tx,
			`
				WITH inserted_google_user AS (
					INSERT INTO accounts.google_users (google_id, given_name, email, profile_picture)
					VALUES ($4, $5, $6, $7)
				)
				INSERT INTO accounts.users (id, google_id, username, profile_picture)
				VALUES ($1, $4, $2, $3)
			`,
			user.ID, user.Username, user.ProfilePicture,
			googleUser.GoogleID, googleUser.GivenName, googleUser.Email, googleUser.ProfilePicture)

		if err != nil {
			return errors.Wrap(err, "cannot insert user and google-user")
		}

		// Upsert a new session for the user.
		sessionID = uuid.New()
		err = database.Insert(ctx, tx, `
			INSERT INTO accounts.sessions (id, user_id, csrf_token)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id) DO UPDATE SET
				id = EXCLUDED.id,
				csrf_token = EXCLUDED.csrf_token
		`,
			sessionID, user.ID, accounts.NewCSRFToken())

		return errors.Wrap(err, "cannot upsert session")
	})

	if err != nil {
		return uuid.Nil, errors.Wrap(err, "transaction failed")
	}

	return sessionID, nil
}

func (db *accountsDatabase) GoogleSignin(ctx context.Context, googleUser *accounts.GoogleUser) (uuid.UUID, error) {

	// Update the google-user.
	err := database.Update(ctx, db.Pool,
		`
			UPDATE accounts.google_users 
			SET given_name = $2, email = $3, profile_picture = $4
			WHERE google_id = $1
		`,
		googleUser.GoogleID, googleUser.GivenName, googleUser.Email, googleUser.ProfilePicture)

	if err != nil {
		return uuid.Nil, errors.Wrap(err, "cannot update google-user by ID")
	}

	// Select google-user's user.
	user, _, err := database.Get[user](ctx, db.Pool,
		`
			SELECT users.id AS "users.id"
			FROM accounts.users
			WHERE users.google_id = $1
		`,
		googleUser.GoogleID)

	if err != nil {
		return uuid.Nil, errors.Wrap(err, "cannot select user by google-ID")
	}

	// Upsert a new session for the user.
	var sessionID = uuid.New()
	err = database.Insert(ctx, db.Pool,
		`
			INSERT INTO accounts.sessions (id, user_id, csrf_token)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id) DO UPDATE SET
				id = EXCLUDED.id,
				csrf_token = EXCLUDED.csrf_token
		`,
		sessionID, user.ID, accounts.NewCSRFToken())

	if err != nil {
		return uuid.Nil, errors.Wrap(err, "cannot upsert session")
	}

	return sessionID, nil
}

package accounts

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	"rollbringer/src/domain/services/accounts"
	"rollbringer/src/repos/database"
)

type spotifyUser struct {
	SpotifyID string `db:"spotify_id"`

	DisplayName    string `db:"display_name"`
	Email          string `db:"email"`
	ProfilePicture string `db:"profile_picture"`
}

func (db *accountsDatabase) SpotifySignup(ctx context.Context, spotifyUser *accounts.SpotifyUser, user *accounts.User) (sessionID uuid.UUID, err error) {
	err = db.Transaction(ctx, func(tx pgx.Tx) error {

		// Insert the user and it's spotify-user.
		err := database.Insert(ctx, tx,
			`
				WITH inserted_spotify_user AS (
					INSERT INTO accounts.spotify_users (spotify_id, display_name, email, profile_picture)
					VALUES ($4, $5, $6, $7)
				)
				INSERT INTO accounts.users (id, spotify_id, username, profile_picture)
				VALUES ($1, $4, $2, $3)
			`,
			user.ID, user.Username, user.ProfilePicture,
			spotifyUser.SpotifyID, spotifyUser.DisplayName, spotifyUser.Email, spotifyUser.ProfilePicture)

		if err != nil {
			return errors.Wrap(err, "cannot insert user and spotify-user")
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

func (db *accountsDatabase) SpotifySignin(ctx context.Context, spotifyUser *accounts.SpotifyUser) (sessionID uuid.UUID, err error) {

	// Update the spotify-user.
	err = database.Update(ctx, db.Pool,
		`
			UPDATE accounts.spotify_users 
			SET display_name = $2, email = $3, profile_picture = $4
			WHERE spotify_id = $1
		`,
		spotifyUser.SpotifyID, spotifyUser.DisplayName, spotifyUser.Email, spotifyUser.ProfilePicture)

	if err != nil {
		return uuid.Nil, errors.Wrap(err, "cannot update spotify-user by ID")
	}

	// Select spotify-user's user.
	user, _, err := database.Get[user](ctx, db.Pool,
		`
			SELECT users.id AS "users.id"
			FROM accounts.users
			WHERE users.spotify_id = $1
		`,
		spotifyUser.SpotifyID)

	if err != nil {
		return uuid.Nil, errors.Wrap(err, "cannot select user by spotify-ID")
	}

	// Upsert a new session for the user.
	sessionID = uuid.New()
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

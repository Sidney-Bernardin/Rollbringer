package accounts

import (
	"context"

	"github.com/pkg/errors"

	"rollbringer/src"
	"rollbringer/src/repositories/database"
	"rollbringer/src/services/accounts"
	"rollbringer/src/services/accounts/models"
)

type spotifyUser struct {
	SpotifyID string `db:"spotify_id"`

	DisplayName    string  `db:"display_name"`
	Email          string  `db:"email"`
	ProfilePicture *string `db:"profile_picture"`
}

/////

const qInsertUserAndSpotifyUser = `
	WITH inserted_spoitfy_user AS (
		INSERT INTO accounts.spotify_users (spotify_id, display_name, email, profile_picture)
		VALUES ($4, $5, $6, $7)
	)
	INSERT INTO accounts.users (id, spotify_id, username, profile_picture)
	VALUES ($1, $4, $2, $3)`

func (db *accountsDatabase) SpotifySignup(ctx context.Context, user *models.User) (sessionID *src.UUID, err error) {
	err = db.Transaction(ctx, func(db *database.Database) error {
		tx := &accountsDatabase{Database: db}

		// Insert the user and user spotify-user.
		err = tx.Insert(ctx, nil, qInsertUserAndSpotifyUser,
			user.UserID, user.Username, user.ProfilePicture,
			user.SpotifyUser.SpotifyID, user.SpotifyUser.DisplayName, user.SpotifyUser.Email, user.SpotifyUser.ProfilePicture)
		if err != nil {
			if errors.Is(err, src.ErrEntityConflict) {
				return &src.ExternalError{
					Type:        accounts.ExternalErrorTypeProviderAlreadyLinked,
					Description: "The Spotify account is already linked with a Rollbringer account.",
				}
			}

			return errors.Wrap(err, "cannot insert spotify-user")
		}

		// Upsert a new session.
		sID := src.NewUUID()
		err = tx.Insert(ctx, nil, qUspertSession,
			sID, user.UserID, models.NewCSRFToken())
		if err != nil {
			return errors.Wrap(err, "cannot upsert session")
		}

		sessionID = &sID
		return nil
	})

	return sessionID, errors.Wrap(err, "transaction failed")
}

/////

const (
	qUpdateSpotifyUserBySpotifyID = `
		UPDATE accounts.spotify_users 
		SET display_name = $2, email = $3, profile_picture = $4
		WHERE spotify_id = $1`

	qSelectUserBySpotifyID = `
		SELECT users.id FROM accounts.users
		WHERE users.spotify_id = $1`
)

func (db *accountsDatabase) SpotifySignin(ctx context.Context, spotifyUser *models.SpotifyUser) (*src.UUID, error) {

	// Update the spotify-user.
	err := db.Update(ctx, nil, qUpdateSpotifyUserBySpotifyID,
		spotifyUser.SpotifyID, spotifyUser.DisplayName, spotifyUser.Email, spotifyUser.ProfilePicture)
	if err != nil {
		if errors.Is(err, src.ErrNoEntitiesEffected) {
			return nil, &src.ExternalError{
				Type:        accounts.ExternalErrorTypeProviderNotLinked,
				Description: "The Spotify account is not linked with a Rollbringer account.",
			}
		}

		return nil, errors.Wrap(err, "cannot update spotify-user by ID")
	}

	// Get the user.
	var u user
	if err = db.SelectOne(ctx, &u, qSelectUserBySpotifyID, spotifyUser.SpotifyID); err != nil {
		return nil, errors.Wrap(err, "cannot select user by spotify-ID")
	}

	// Upsert a new session.
	sessionID := src.NewUUID()
	err = db.Insert(ctx, nil, qUspertSession,
		sessionID, u.ID, models.NewCSRFToken())
	if err != nil {
		return nil, errors.Wrap(err, "cannot upsert session")
	}

	return &sessionID, nil
}

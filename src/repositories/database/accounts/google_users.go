package accounts

import (
	"context"

	"github.com/pkg/errors"

	"rollbringer/src"
	"rollbringer/src/repositories/database"
	"rollbringer/src/services/accounts"
	"rollbringer/src/services/accounts/models"
)

type googleUser struct {
	GoogleID string `db:"google_id"`

	GivenName      string `db:"given_name"`
	Email          string `db:"email"`
	ProfilePicture string `db:"profile_picture"`
}

/////

const qInsertUserAndGoogleUser = `
	WITH inserted_google_user AS (
		INSERT INTO accounts.google_users (google_id, given_name, email, profile_picture)
		VALUES ($4, $5, $6, $7)
	)
	INSERT INTO accounts.users (id, google_id, username, profile_picture)
	VALUES ($1, $4, $2, $3)`

func (db *accountsDatabase) GoogleSignup(ctx context.Context, user *models.User) (sessionID *src.UUID, err error) {
	err = db.Transaction(ctx, func(db *database.Database) error {
		tx := &accountsDatabase{Database: db}

		// Insert the user and user google-user.
		err = tx.Insert(ctx, nil, qInsertUserAndGoogleUser,
			user.UserID, user.Username, user.ProfilePicture,
			user.GoogleUser.GoogleID, user.GoogleUser.GivenName, user.GoogleUser.Email, user.GoogleUser.ProfilePicture)
		if err != nil {
			if errors.Is(err, src.ErrEntityConflict) {
				return &src.ExternalError{
					Type:        accounts.ExternalErrorTypeProviderAlreadyLinked,
					Description: "The Google account is already linked with a Rollbringer account.",
				}
			}

			return errors.Wrap(err, "cannot insert google-user")
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
	qUpdateGoogleUserByGoogleID = `
		UPDATE accounts.google_users 
		SET given_name = $2, email = $3, profile_picture = $4
		WHERE google_id = $1`

	qSelectUserByGoogleID = `
		SELECT users.id FROM accounts.users
		WHERE users.google_id = $1`
)

func (db *accountsDatabase) GoogleSignin(ctx context.Context, googleUser *models.GoogleUser) (*src.UUID, error) {

	// Update the google-user.
	err := db.Update(ctx, nil, qUpdateGoogleUserByGoogleID,
		googleUser.GoogleID, googleUser.GivenName, googleUser.Email, googleUser.ProfilePicture)
	if err != nil {
		if errors.Is(err, src.ErrNoEntitiesEffected) {
			return nil, &src.ExternalError{
				Type:        accounts.ExternalErrorTypeProviderNotLinked,
				Description: "The Google account is not linked with a Rollbringer account.",
			}
		}

		return nil, errors.Wrap(err, "cannot update google-user by ID")
	}

	// Get the user.
	var u user
	if err = db.SelectOne(ctx, &u, qSelectUserByGoogleID, googleUser.GoogleID); err != nil {
		return nil, errors.Wrap(err, "cannot select user by google-ID")
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

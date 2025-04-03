package accounts

import (
	"context"
	"embed"

	"github.com/pkg/errors"

	"rollbringer/src"
	"rollbringer/src/domain"
	"rollbringer/src/domain/accounts"
	"rollbringer/src/repositories/database"
)

//go:embed migrations/*.sql
var migrations embed.FS

type accountsDatabase struct {
	*database.Database
}

func NewDatabase(config *src.Config) (accounts.Database, error) {
	database, err := database.NewDatabase(config.PostgresAccountsURL, &migrations)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create database")
	}

	return &accountsDatabase{
		Database: database,
	}, nil
}

func (db *accountsDatabase) GoogleSignup(ctx context.Context, user *accounts.User) (sessionID domain.UUID, err error) {
	err = db.Transaction(ctx, func(db *database.Database) error {
		tx := &accountsDatabase{Database: db}

		// Insert the google-user.
		err = tx.CRUDInsert(ctx, nil, qGoogleUserInsert,
			user.GoogleUser.GoogleID, user.GoogleUser.GivenName, user.GoogleUser.Email, user.GoogleUser.ProfilePicture)
		if err != nil {
			if errors.Is(err, domain.ErrEntityConflict) {
				return &src.ExternalError{
					Type:        accounts.ExternalErrorTypeProviderAlreadyLinked,
					Description: "The Google account is already linked with a Rollbringer account.",
				}
			}

			return errors.Wrap(err, "cannot insert google-user")
		}

		// Insert the user.
		err = tx.CRUDInsert(ctx, nil, qInsertUser,
			user.UserID, user.GoogleID, user.SpotifyID, user.Username, user.ProfilePicture)
		if err != nil {
			return errors.Wrap(err, "cannot insert user")
		}

		// Upsert a new session.
		sessionID = domain.NewUUID()
		err = tx.CRUDInsert(ctx, nil, qSessionUpsert,
			sessionID, user.UserID, accounts.NewCSRFToken())
		return errors.Wrap(err, "cannot insert session")
	})

	return sessionID, errors.Wrap(err, "transaction failed")
}

func (db *accountsDatabase) GoogleSignin(ctx context.Context, googleUser *accounts.GoogleUser) (sessionID domain.UUID, err error) {

	// Update the google-user.
	err = db.CRUDUpdate(ctx, nil, qGoogleUserUpdateByID,
		googleUser.GoogleID, googleUser.GivenName, googleUser.Email, googleUser.ProfilePicture)
	if err != nil {
		if errors.Is(err, domain.ErrNoEntitiesEffected) {
			return domain.UUID{}, &src.ExternalError{
				Type:        accounts.ExternalErrorTypeProviderNotLinked,
				Description: "The Google account is not linked with a Rollbringer account.",
			}
		}

		return domain.UUID{}, errors.Wrap(err, "cannot update google-user by ID")
	}

	// Get the user.
	var userID domain.UUID
	err = db.queryUser(ctx, db.CRUDGet, &userID, qUserSelectByGoogleID, googleUser.GoogleID)
	if err != nil {
		return domain.UUID{}, errors.Wrap(err, "cannot get user by google-id")
	}

	// Upsert a new session.
	sessionID = domain.NewUUID()
	err = db.CRUDInsert(ctx, nil, qSessionUpsert,
		sessionID, userID, accounts.NewCSRFToken())
	return sessionID, errors.Wrap(err, "cannot insert session")
}

func (db *accountsDatabase) SpotifySignup(ctx context.Context, user *accounts.User) (sessionID domain.UUID, err error) {
	err = db.Transaction(ctx, func(db *database.Database) error {
		tx := &accountsDatabase{Database: db}

		// Insert the spotify-user.
		err = tx.CRUDInsert(ctx, nil, qSpotifyUserInsert,
			user.SpotifyUser.SpotifyID, user.SpotifyUser.DisplayName, user.SpotifyUser.Email, user.SpotifyUser.ProfilePicture)
		if err != nil {
			if errors.Is(err, domain.ErrEntityConflict) {
				return &src.ExternalError{
					Type:        accounts.ExternalErrorTypeProviderAlreadyLinked,
					Description: "The Google account is already linked with a Rollbringer account.",
				}
			}

			return errors.Wrap(err, "cannot insert spotify-user")
		}

		// Insert the user.
		err := tx.CRUDInsert(ctx, nil, qInsertUser,
			user.UserID, user.GoogleID, user.SpotifyID, user.Username, user.ProfilePicture)
		if err != nil {
			return errors.Wrap(err, "cannot insert user")
		}

		// Upsert a new session.
		sessionID = domain.NewUUID()
		err = tx.CRUDInsert(ctx, nil, qSessionUpsert,
			sessionID, user.UserID, accounts.NewCSRFToken())
		return errors.Wrap(err, "cannot insert session")
	})

	return sessionID, errors.Wrap(err, "transaction failed")
}

func (db *accountsDatabase) SpotifySignin(ctx context.Context, spotifyUser *accounts.SpotifyUser) (sessionID domain.UUID, err error) {

	// Update the spotify-user.
	err = db.CRUDUpdate(ctx, nil, qSpotifyUserUpdateByID,
		spotifyUser.SpotifyID, spotifyUser.DisplayName, spotifyUser.Email, spotifyUser.ProfilePicture)
	if err != nil {
		if errors.Is(err, domain.ErrNoEntitiesEffected) {
			return domain.UUID{}, &src.ExternalError{
				Type:        accounts.ExternalErrorTypeProviderNotLinked,
				Description: "The Spotify account is not linked with a Rollbringer account.",
			}
		}

		return domain.UUID{}, errors.Wrap(err, "cannot update spotify-user by ID")
	}

	// Get the user.
	var userID domain.UUID
	err = db.queryUser(ctx, db.CRUDGet, &userID, qUserSelectBySpotifyID, spotifyUser.SpotifyID)
	if err != nil {
		return domain.UUID{}, errors.Wrap(err, "cannot get user by spotify-id")
	}

	// Upsert a new session.
	sessionID = domain.NewUUID()
	err = db.CRUDInsert(ctx, nil, qSessionUpsert,
		sessionID, userID, accounts.NewCSRFToken())
	return sessionID, errors.Wrap(err, "cannot insert session")
}

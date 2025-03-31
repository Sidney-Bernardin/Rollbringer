package accounts

import (
	"context"
	"embed"

	"rollbringer/src"
	"rollbringer/src/domain/accounts"
	"rollbringer/src/repositories/database"

	"github.com/google/uuid"
	"github.com/pkg/errors"
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

func (db *accountsDatabase) GoogleSignup(ctx context.Context, user *accounts.User) (uuid.UUID, error) {
	if err := db.queryUser(ctx, db.CRUDInsert, nil, qInsertUser, user.UserID, user.Username); err != nil {
		return uuid.Nil, errors.Wrap(err, "cannot insert user")
	}

	err := db.CRUDInsert(ctx, nil, qGoogleUserInsert,
		user.GoogleUser.GoogleID, user.GoogleUser.GivenName, user.GoogleUser.Email, user.GoogleUser.ProfilePicture)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "cannot insert google-user")
	}

	var sessionID uuid.UUID
	err = db.querySession(ctx, db.CRUDInsert, sessionID, qSessionInsert,
		uuid.New(), user.UserID, accounts.NewCSRFToken())
	return sessionID, errors.Wrap(err, "cannot insert session")
}

func (db *accountsDatabase) GoogleSignin(ctx context.Context, googleUser *accounts.GoogleUser) (uuid.UUID, error) {
	err := db.CRUDUpdate(ctx, nil, qGoogleUserUpdateByID,
		googleUser.GoogleID, googleUser.GivenName, googleUser.Email, googleUser.ProfilePicture)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "cannot update google-user by ID")
	}

	var userID uuid.UUID
	err = db.queryUser(ctx, db.CRUDGet, userID, qUserSelectByGoogleID, googleUser.GoogleID)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "cannot get user by google-id")
	}

	var sessionID uuid.UUID
	err = db.querySession(ctx, db.CRUDInsert, sessionID, qSessionInsert,
		uuid.New(), userID, accounts.NewCSRFToken())
	return sessionID, errors.Wrap(err, "cannot insert session")
}

func (db *accountsDatabase) SpotifySignup(ctx context.Context, user *accounts.User) (uuid.UUID, error) {
	if err := db.queryUser(ctx, db.CRUDInsert, nil, qInsertUser, user.UserID, user.Username); err != nil {
		return uuid.Nil, errors.Wrap(err, "cannot insert user")
	}

	err := db.CRUDInsert(ctx, nil, qSpotifyUserInsert,
		user.GoogleUser.GoogleID, user.GoogleUser.GivenName, user.GoogleUser.Email, user.GoogleUser.ProfilePicture)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "cannot insert spotify-user")
	}

	var sessionID uuid.UUID
	err = db.querySession(ctx, db.CRUDInsert, sessionID, qSessionInsert,
		uuid.New(), user.UserID, accounts.NewCSRFToken())
	return sessionID, errors.Wrap(err, "cannot insert session")
}

func (db *accountsDatabase) SpotifySignin(ctx context.Context, spotifyUser *accounts.SpotifyUser) (uuid.UUID, error) {
	err := db.CRUDUpdate(ctx, nil, qSpotifyUserUpdateByID,
		spotifyUser.SpotifyID, spotifyUser.DisplayName, spotifyUser.Email, spotifyUser.ProfilePicture)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "cannot update spotify-user by ID")
	}

	var userID uuid.UUID
	err = db.queryUser(ctx, db.CRUDGet, userID, qUserSelectByGoogleID, spotifyUser.SpotifyID)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "cannot get user by google-id")
	}

	var sessionID uuid.UUID
	err = db.querySession(ctx, db.CRUDInsert, sessionID, qSessionInsert,
		uuid.New(), userID, accounts.NewCSRFToken())
	return sessionID, errors.Wrap(err, "cannot insert session")
}

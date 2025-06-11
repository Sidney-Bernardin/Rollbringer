package service

import (
	"context"
	"fmt"

	"github.com/Sidney-Bernardin/Rollbringer/server"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/google"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/sql"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/sql/queries"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (svc *Service) BasicSignup(ctx context.Context, username, password string) (sessionID server.UUID, err error) {

	if len(username) < 3 {
		return sessionID, &server.UserError{
			Type:    server.UserErrorTypePasswordInvalid,
			Message: "Password must at least 3 characters long.",
		}
	}

	if len(username) < 3 || 32 < len(username) {
		return sessionID, &server.UserError{
			Type:    server.UserErrorTypeUsernameInvalid,
			Message: "Username must be between 3 and 32 characters long.",
		}
	}

	var (
		userID       = server.NewRandomUUID()
		passwordSalt = server.CreateRandomString()
		passwordHash = []byte(password + passwordSalt)
	)

	passwordHash, err = bcrypt.GenerateFromPassword(passwordHash, 12)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return sessionID, &server.UserError{
				Type:    server.UserErrorTypePasswordInvalid,
				Message: "Password too long",
			}
		}

		return sessionID, errors.Wrap(err, "cannot hash password")
	}

	affected, err := svc.SQL.InsertUser(ctx, &queries.InsertUserParams{
		ID:           userID,
		Username:     username,
		PasswordHash: passwordHash,
		PasswordSalt: &passwordSalt,
	})

	if err != nil {
		return sessionID, errors.Wrap(err, "cannot create user")
	} else if affected < 1 {
		return sessionID, &server.UserError{
			Type:    server.UserErrorTypeUsernameTaken,
			Message: "That username is already being used by another Rollbringer user.",
		}
	}

	sessionID, err = svc.Cache.SetSession(ctx, userID)
	return sessionID, errors.Wrap(err, "cannot put session")
}

func (svc *Service) BasicSignin(ctx context.Context, username, password string) (sessionID server.UUID, err error) {

	user, err := svc.SQL.SelectUserWithPassword(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sessionID, &server.UserError{Type: server.UserErrorTypeUnauthorized}
		}

		return sessionID, errors.Wrap(err, "cannot get user")
	}

	if err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password+*user.PasswordSalt)); err != nil {
		return sessionID, &server.UserError{Type: server.UserErrorTypeUnauthorized}
	}

	sessionID, err = svc.Cache.SetSession(ctx, user.ID)
	return sessionID, errors.Wrap(err, "cannot put session")
}

func (svc *Service) GoogleSignup(ctx context.Context, googleUser *google.GoogleUser) (sessionID server.UUID, err error) {
	userID := server.NewRandomUUID()

	err = svc.SQL.Transaction(ctx, func(tx *sql.SQL) error {

		affected, err := tx.InsertGoogleUser(ctx, &queries.InsertGoogleUserParams{
			GoogleID:  googleUser.Subject,
			GivenName: googleUser.GivenName,
			Email:     googleUser.Email,
		})

		if err != nil {
			return errors.Wrap(err, "cannot create google-user")
		} else if affected < 1 {
			return &server.UserError{
				Type:    server.UserErrorTypeGoogleUserAlreadyExists,
				Message: "That Google account is being used by another Rollbringer user.",
			}
		}

		_, err = tx.InsertUser(ctx, &queries.InsertUserParams{
			ID:             userID,
			GoogleID:       &googleUser.Subject,
			Username:       fmt.Sprintf("%s %s", googleUser.GivenName, googleUser.Subject),
			ProfilePicture: googleUser.Picture,
		})

		return errors.Wrap(err, "cannot create user")
	})

	if err != nil {
		return sessionID, errors.WithStack(err)
	}

	sessionID, err = svc.Cache.SetSession(ctx, userID)
	return sessionID, errors.Wrap(err, "cannot put session")
}

func (svc *Service) GoogleSignin(ctx context.Context, googleUser *google.GoogleUser) (sessionID server.UUID, err error) {
	var userID server.UUID

	err = svc.SQL.Transaction(ctx, func(tx *sql.SQL) error {

		affected, err := tx.UpdateGoogleUser(ctx, &queries.UpdateGoogleUserParams{
			GoogleID:  googleUser.Subject,
			GivenName: googleUser.GivenName,
			Email:     googleUser.Email,
		})

		if err != nil {
			return errors.Wrap(err, "cannot update google-user")
		} else if affected <= 0 {
			return &server.UserError{
				Type:    server.UserErrorTypeGoogleUserNotExists,
				Message: "That Google account isn't being used by another Rollbringer user.",
			}
		}

		userID, err = tx.SelectUserID(ctx, &googleUser.Subject)
		return errors.Wrap(err, "cannot get user")
	})

	if err != nil {
		return sessionID, errors.WithStack(err)
	}

	sessionID, err = svc.Cache.SetSession(ctx, userID)
	return sessionID, errors.Wrap(err, "cannot put session")
}

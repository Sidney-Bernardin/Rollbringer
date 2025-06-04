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
		return sessionID, server.NewUserError(server.UserErrorTypePasswordInvalid, "Password must at least 3 characters long.", nil)
	}

	if len(username) < 3 || 32 < len(username) {
		return sessionID, server.NewUserError(server.UserErrorTypeUsernameInvalid, "Username must be between 3 and 32 characters long.", nil)
	}

	var (
		userID       = server.NewRandomUUID()
		passwordSalt = server.CreateRandomString()
		passwordHash = []byte(password + passwordSalt)
	)

	passwordHash, err = bcrypt.GenerateFromPassword(passwordHash, 12)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return sessionID, server.NewUserError(server.UserErrorTypePasswordInvalid, "Password too long", nil)
		}

		return sessionID, errors.Wrap(err, "cannot hash password")
	}

	affected, err := svc.SQL.CreateUser(ctx, &queries.CreateUserParams{
		ID:           userID,
		Username:     username,
		PasswordHash: passwordHash,
		PasswordSalt: &passwordSalt,
	})

	if err != nil {
		return sessionID, errors.Wrap(err, "cannot create user")
	} else if affected < 1 {
		return sessionID, server.NewUserError(server.UserErrorTypeUsernameTaken, "", nil)
	}

	sessionID, err = svc.Nats.PutSession(ctx, userID)
	return sessionID, errors.Wrap(err, "cannot put session")
}

func (svc *Service) BasicSignin(ctx context.Context, username, password string) (sessionID server.UUID, err error) {

	user, err := svc.SQL.GetUserWithPassword(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sessionID, server.NewUserError(server.UserErrorTypeUnauthorized, "", nil)
		}

		return sessionID, errors.Wrap(err, "cannot get user")
	}

	if err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password+*user.PasswordSalt)); err != nil {
		return sessionID, server.NewUserError(server.UserErrorTypeUnauthorized, "", nil)
	}

	sessionID, err = svc.Nats.PutSession(ctx, user.ID)
	return sessionID, errors.Wrap(err, "cannot put session")
}

func (svc *Service) GoogleSignup(ctx context.Context, googleUser *google.GoogleUser) (sessionID server.UUID, err error) {
	userID := server.NewRandomUUID()

	err = svc.SQL.Transaction(ctx, func(tx *sql.SQL) error {

		affected, err := tx.CreateGoogleUser(ctx, &queries.CreateGoogleUserParams{
			GoogleID:  googleUser.Subject,
			GivenName: googleUser.GivenName,
			Email:     googleUser.Email,
		})

		if err != nil {
			return errors.Wrap(err, "cannot create google-user")
		} else if affected < 1 {
			return server.NewUserError(server.UserErrorTypeGoogleUserAlreadyExists,
				"That Google account is being used by another Rollbringer user.", nil)
		}

		_, err = tx.CreateUser(ctx, &queries.CreateUserParams{
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

	sessionID, err = svc.Nats.PutSession(ctx, userID)
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
			return server.NewUserError(server.UserErrorTypeGoogleUserNotExists,
				"That Google account isn't being used by another Rollbringer user.", nil)
		}

		userID, err = tx.GetUserID(ctx, &googleUser.Subject)
		return errors.Wrap(err, "cannot get user")
	})

	if err != nil {
		return sessionID, errors.WithStack(err)
	}

	sessionID, err = svc.Nats.PutSession(ctx, userID)
	return sessionID, errors.Wrap(err, "cannot put session")
}

package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/Sidney-Bernardin/Rollbringer/server"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/google"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/sql"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/sql/queries"
	"github.com/pkg/errors"
)

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
		} else if affected <= 0 {
			return server.NewUserError(server.UserErrorTypeGoogleUserAlreadyExists,
				"That Google account is being used by another Rollbringer user.", nil)
		}

		err = tx.CreateUser(ctx, &queries.CreateUserParams{
			ID:             userID,
			GoogleID:       &googleUser.Subject,
			Username:       googleUser.GivenName,
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

func (svc *Service) GetUser(ctx context.Context, userID server.UUID) (*queries.User, error) {

	user, err := svc.Nats.GetUser(ctx, userID)
	if err != nil || user != nil {
		return user, errors.Wrap(err, "cannot get user from Nats")
	}

	user, err = svc.SQL.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, server.NewUserError(server.UserErrorTypeUserNotFound, "", nil)
		}

		return nil, errors.Wrap(err, "cannot get user from SQL")
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := svc.Nats.PutUser(ctx, user); err != nil {
			svc.Log.Log(ctx, slog.LevelWarn, "Nats put user", "err", err.Error())
		}
	}()

	return user, nil
}

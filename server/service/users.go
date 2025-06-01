package service

import (
	"context"

	"github.com/Sidney-Bernardin/Rollbringer/server"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/google"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/sql"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/sql/queries"
	"github.com/pkg/errors"
)

func (svc *Service) GoogleSignup(ctx context.Context, googleUser *google.GoogleUser) (sessionID server.UUID, err error) {
	userID := server.NewRandomUUID()

	err = svc.sql.Transaction(ctx, func(tx *sql.SQL) error {

		affected, err := tx.InsertGoogleUser(ctx, &queries.InsertGoogleUserParams{
			GoogleID:  googleUser.ID,
			GivenName: googleUser.GivenName,
			Email:     googleUser.Email,
		})

		if err != nil {
			return errors.Wrap(err, "cannot insert google user")
		} else if affected == 0 {
			return &server.UserError{
				Type:    server.UserErrorTypeGoogleUserAlreadyExists,
				Message: "A user with that Google account already exists.",
			}
		}

		err = tx.InsertUser(ctx, &queries.InsertUserParams{
			ID:             userID,
			GoogleID:       &googleUser.ID,
			Username:       googleUser.GivenName,
			ProfilePicture: googleUser.Picture,
		})

		return errors.Wrap(err, "cannot insert user")
	})

	if err != nil {
		return sessionID, errors.Wrap(err, "cannot insert user info")
	}

	sessionID, err = svc.nats.PutSession(ctx, userID)
	return sessionID, errors.Wrap(err, "cannot put session")
}

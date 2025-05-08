package accounts

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"rollbringer/src/domain/services/accounts"
	"rollbringer/src/repos/database"
)

type user struct {
	ID             uuid.UUID `db:"users.id"`
	GoogleID       *string   `db:"users.google_id"`
	SpotifyID      *string   `db:"users.spotify_id"`
	Username       string    `db:"users.username"`
	ProfilePicture string    `db:"users.profile_picture"`
}

func (u user) Model() *accounts.User {
	return &accounts.User{
		ID:             u.ID,
		GoogleID:       u.GoogleID,
		SpotifyID:      u.SpotifyID,
		Username:       accounts.Username(u.Username),
		ProfilePicture: u.ProfilePicture,
	}
}

func (db *accountsDatabase) GetUserByUserID(ctx context.Context, userID uuid.UUID) (*accounts.User, error) {
	_, model, err := database.Get[user](ctx, db.Pool,
		`
			SELECT
				users.id AS "users.id",
				users.username AS "users.username",
				users.profile_picture AS "users.profile_picture"
			FROM accounts.users
			WHERE users.id = $1
		`,
		userID)

	return model, errors.Wrap(err, "cannot select user by user-ID")
}

func (db *accountsDatabase) GetUsersByUserIDs(ctx context.Context, userIDs ...uuid.UUID) ([]*accounts.User, error) {
	if len(userIDs) < 1 {
		return []*accounts.User{}, nil
	}

	_, models, err := database.Gets[user](ctx, db.Pool,
		`
			SELECT
				users.id AS "users.id",
				users.username AS "users.username",
				users.profile_picture AS "users.profile_picture"
			FROM accounts.users
			WHERE users.id = ANY($1)
		`,
		userIDs)

	return models, errors.Wrap(err, "cannot select users by user-IDs")
}

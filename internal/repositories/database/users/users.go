package users

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"rollbringer/internal"
	"rollbringer/internal/repositories/database"
)

func userColumns(views map[string]internal.UserView) (columns string) {
	switch views["user"] {
	case internal.UserViewUserAll:
		columns += `users.*`
	default:
		columns += `users.*`
	}

	return columns
}

func (db *usersSchema) UserInsert(ctx context.Context, user *internal.User) error {
	err := db.TX.QueryRowxContext(ctx,
		`INSERT INTO users.users (id, google_id, username)
			VALUES ($1, $2, $3)
		ON CONFLICT (google_id)
			DO UPDATE SET google_id = EXCLUDED.google_id
		RETURNING id`,
		uuid.New(), user.GoogleID, user.Username,
	).Scan(&user.ID)

	return errors.Wrap(err, "cannot insert user")
}

func (db *usersSchema) UserGet(ctx context.Context, userID uuid.UUID, views map[string]internal.UserView) (*internal.User, error) {

	var user database.User
	query := fmt.Sprintf(`SELECT %s FROM users.users WHERE id = $1`, userColumns(views))
	if err := sqlx.GetContext(ctx, db.TX, &user, query, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, internal.NewProblemDetail(ctx, internal.PDOpts{
				Type:   internal.PDTypeUserNotFound,
				Detail: "Can't find a user with the given user_id.",
				Extra: map[string]any{
					"user_id": userID,
				},
			})
		}

		return nil, errors.Wrap(err, "cannot select user")
	}

	return user.Internalized(), nil
}

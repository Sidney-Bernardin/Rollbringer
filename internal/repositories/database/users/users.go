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
	userView, ok := views["user"]
	if !ok {
		userView = views["users"]
	}

	switch userView {
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

func (db *usersSchema) UsersGetByGame(ctx context.Context, gameID uuid.UUID, views map[string]internal.UserView) ([]*internal.User, error) {
	query := fmt.Sprintf(
		`SELECT %s FROM users.users
		WHERE EXISTS (
			SELECT * FROM game_users WHERE game_users.game_id = $1 AND game_users.user_id = users.id
		)`,
		userColumns(views),
	)

	var users []*database.User
	if err := sqlx.SelectContext(ctx, db.TX, &users, query, gameID); err != nil {
		return nil, errors.Wrap(err, "cannot select users")
	}

	// Internalize each user.
	ret := make([]*internal.User, len(users))
	for i, m := range users {
		ret[i] = m.Internalized()
	}

	return ret, nil
}

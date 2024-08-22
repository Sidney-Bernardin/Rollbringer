package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"rollbringer/internal"
)

var userViewColumns = map[internal.UserView]string{
	internal.UserViewAll: "users.id, users.google_id, users.username",
}

type dbUser struct {
	ID uuid.UUID `db:"id"`

	GoogleID *string `db:"google_id"`
	Username string  `db:"username"`
}

func (user *dbUser) internalized() *internal.User {
	if user != nil {
		return &internal.User{
			ID:       user.ID,
			GoogleID: user.GoogleID,
			Username: user.Username,
		}
	}
	return nil
}

func (db *UsersDatabase) UserInsert(ctx context.Context, user *internal.User) error {
	err := db.TX.QueryRowxContext(ctx,
		`INSERT INTO users (id, google_id, username)
			VALUES ($1, $2, $3)
		ON CONFLICT (google_id)
			DO UPDATE SET google_id = EXCLUDED.google_id
		RETURNING id`,
		uuid.New(), user.GoogleID, user.Username,
	).Scan(&user.ID)

	return errors.Wrap(err, "cannot insert user")
}

func (db *UsersDatabase) UserGet(ctx context.Context, userID uuid.UUID, view internal.UserView) (*internal.User, error) {
	columns, ok := userViewColumns[view]
	if !ok {
		return nil, fmt.Errorf("bad user view %d", view)
	}
	query := fmt.Sprintf(`SELECT %s FROM users WHERE id = $1`, columns)

	var user dbUser
	if err := sqlx.GetContext(ctx, db.TX, &user, query, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, internal.NewProblemDetail(ctx, &internal.PDOptions{
				Type:   internal.PDTypeUserNotFound,
				Detail: "Can't find a user with the given user-ID.",
				Extra: map[string]any{
					"user_id": userID,
				},
			})
		}

		return nil, errors.Wrap(err, "cannot select user")
	}

	return user.internalized(), nil
}

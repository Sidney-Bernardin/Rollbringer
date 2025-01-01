package users

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"rollbringer/internal"
	"rollbringer/internal/repositories/database"
)

func (db *usersSchema) UserInsert(ctx context.Context, user *internal.User) error {
	err := db.TX.QueryRowxContext(ctx,
		`INSERT INTO users.users (id, username, google_id, google_picture)
			VALUES ($1, $2, $3, $4)
		ON CONFLICT (google_id)
			DO UPDATE SET google_id = EXCLUDED.google_id
		RETURNING id`,
		uuid.New(), user.Username, user.GoogleID, user.GooglePicture,
	).Scan(&user.ID)

	return errors.Wrap(err, "cannot insert user")
}

func (db *usersSchema) UsersGetByGame(ctx context.Context, gameID uuid.UUID) ([]*internal.User, error) {
	query := `
		SELECT users.* FROM users.users
		WHERE EXISTS (
			SELECT * FROM game_users WHERE game_users.game_id = $1 AND game_users.user_id = users.id
		)
	`

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

package sql

import (
	"context"

	"github.com/Sidney-Bernardin/Rollbringer/internal/domain"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func (repo *repository) GoogleSignup(ctx context.Context, user *domain.User, googleUser *domain.GoogleUser) error {
	return errors.WithStack(repo.transaction(ctx, func(tx *repository) error {

		_, err := sqlx.NamedExecContext(ctx, repo.tx,
			`
				INSERT INTO rollbringer.google_users (google_id, given_name, email)
				VALUES (:google_id, :given_name, :email)
			`,
			googleUser)

		if err != nil {
			return errors.Wrap(parseInsertErr[domain.GoogleUser](err), "cannot insert google-user")
		}

		_, err = sqlx.NamedExecContext(ctx, repo.tx,
			`
				INSERT INTO rollbringer.users (id, google_id, username, profile_picture)
				VALUES (:id, :google_id, :username, :profile_picture)
			`,
			user)

		return errors.Wrap(parseInsertErr[domain.User](err), "cannot insert user")
	}))
}

func (repo *repository) GoogleSignin(ctx context.Context, googleUser *domain.GoogleUser) (userID domain.UUID, err error) {
	result, err := sqlx.NamedExecContext(ctx, repo.tx,
		`
			UPDATE rollbringer.google_users
			SET given_name = :given_name, email = :email
			WHERE google_id = :google_id
		`,
		googleUser)

	if err != nil {
		return userID, errors.Wrap(parseUpdateErr[domain.GoogleUser](result, err), "cannot update google-user")
	}

	var user *domain.User
	err = sqlx.GetContext(ctx, repo.tx, &user,
		`
			SELECT id FROM rollbringer.users
			WHERE google_id = $1
		`,
		googleUser.GoogleID)

	if err != nil {
		return userID, errors.Wrap(parseGetErr[domain.User](err), "cannot select user")
	}

	return user.ID, nil
}

func (repo *repository) GetUser(ctx context.Context, userID domain.UUID) (*domain.User, error) {
	var user domain.User

	err := sqlx.GetContext(ctx, repo.tx, &user,
		`
			SELECT id, created_at, updated_at, username, profile_picture
			FROM rollbringer.users
			WHERE id = $1
		`,
		userID)

	return &user, errors.Wrap(parseGetErr[domain.User](err), "cannot select user")
}

package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
)

var userViewColumns = map[domain.UserView]string{
	domain.UserViewAll: "users.id, users.google_id, users.username",
}

type userModel struct {
	ID uuid.UUID `db:"id"`

	GoogleID *string `db:"google_id"`
	Username string  `db:"username"`
}

func (user *userModel) domain() *domain.User {
	if user != nil {
		return &domain.User{
			ID:       user.ID,
			GoogleID: user.GoogleID,
			Username: user.Username,
		}
	}
	return nil
}

func (db *Database) InsertUser(ctx context.Context, user *domain.User) error {

	model := userModel{
		ID:       uuid.New(),
		GoogleID: user.GoogleID,
		Username: user.Username,
	}

	// Insert the user.
	err := sqlx.GetContext(ctx, db.tx, &model,
		`INSERT INTO users (id, google_id, username)
			VALUES ($1, $2, $3)
		ON CONFLICT (google_id)
			DO UPDATE SET google_id = EXCLUDED.google_id
		RETURNING id`,
		model.ID, model.GoogleID, model.Username,
	)

	if err != nil {
		return errors.Wrap(err, "cannot insert user")
	}

	*user = *model.domain()
	return nil
}

func (db *Database) GetJoinedUsersForGame(ctx context.Context, gameID uuid.UUID, view domain.UserView) ([]*domain.User, error) {

	// Build a query to select joined users with the game.
	query := fmt.Sprintf(
		`SELECT %s FROM users
		WHERE EXISTS (
			SELECT * FROM game_joined_users WHERE game_joined_users.game_id = $1 AND game_joined_users.user_id = users.id
		)`,
		userViewColumns[view],
	)

	// Execute the query.
	var models []*userModel
	if err := sqlx.SelectContext(ctx, db.tx, &models, query, gameID); err != nil {
		return nil, errors.Wrap(err, "cannot select games")
	}

	// Convert each model to domain.User.
	ret := make([]*domain.User, len(models))
	for i, m := range models {
		ret[i] = m.domain()
	}

	return ret, nil
}

func (db *Database) GetUser(ctx context.Context, userID uuid.UUID, view domain.UserView) (*domain.User, error) {

	// Build a query to select a user with the user-ID.
	query := fmt.Sprintf(
		`SELECT %s FROM users WHERE id = $1`,
		userViewColumns[view],
	)

	// Execute the query.
	var model userModel
	if err := sqlx.GetContext(ctx, db.tx, &model, query, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, &domain.NormalError{
				Type:   domain.NETypeUserNotFound,
				Detail: fmt.Sprintf("Cannot find a user with the user-ID"),
			}
		}

		return nil, errors.Wrap(err, "cannot select user")
	}

	return model.domain(), nil
}

func (db *Database) GetUserByGoogleID(ctx context.Context, googleID string, view domain.UserView) (*domain.User, error) {
	return &domain.User{
		ID:          [16]byte{},
		GoogleID:    &googleID,
		Username:    "",
		PDFs:        []*domain.PDF{},
		HostedGames: []*domain.Game{},
		JoinedGames: []*domain.Game{},
	}, nil
}

package database

import (
	"context"
	"fmt"

	"rollbringer/pkg/domain"
)

/////

const qUserInsert = ` 
WITH inserted_user AS (
	INSERT INTO accounts.users (google_id, spotify_id, username, profile_picture)
	VALUES ($1, $2, $3, $4)
	RETURNING *
)
SELECT * FROM inserted_user`

func (repo *accountsDatabaseRepository) UserInsert(ctx context.Context, user *domain.User) error {
	err := repo.Insert(ctx, user, qUserInsert,
		user.GoogleID, user.SpotifyID, user.Username, user.ProfilePicture)
	return domain.Wrap(err, "cannot insert user", nil)
}

/////

const qUserGet = ` 
SELECT * FROM accounts.users WHERE users.%s = $1`

func (repo *accountsDatabaseRepository) UserGet(ctx context.Context, key string, value any) (*domain.User, error) {
	user := &domain.User{}
	if err := repo.Get(ctx, user, fmt.Sprintf(qUserGet, key), value); err != nil {
		return nil, domain.Wrap(err, "cannot select user", nil)
	}
	return user, nil
}

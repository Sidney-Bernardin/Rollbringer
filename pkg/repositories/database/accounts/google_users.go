package database

import (
	"context"
	"fmt"
	"rollbringer/pkg/domain"
)

/////

const qGoogleUserInsert = ` 
INSERT INTO accounts.google_users (google_id, given_name, email, profile_picture)
VALUES ($1, $2, $3, $4)
RETURNING *`

func (repo *accountsDatabaseRepository) GoogleUserInsert(ctx context.Context, googleUser *domain.GoogleUser) error {
	err := repo.Insert(ctx, googleUser, qGoogleUserInsert,
		googleUser.GoogleID, googleUser.GivenName, googleUser.Email, googleUser.ProfilePicture)
	return domain.Wrap(err, "cannot insert google-user", nil)
}

/////

const qGoogleUserUpdate = ` 
UPDATE accounts.google_users
SET %s
WHERE %s = $1`

func (repo *accountsDatabaseRepository) GoogleUserUpdate(ctx context.Context, key string, value any, updates map[string]any) error {
	err := repo.Update(ctx, updates, fmt.Sprintf(qGoogleUserUpdate, "%s", key), value)
	return domain.Wrap(err, "cannot update google-user", nil)
}

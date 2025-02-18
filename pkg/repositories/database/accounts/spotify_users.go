package database

import (
	"context"
	"fmt"
	"rollbringer/pkg/domain"
)

const qSpotifyUserInsert = ` 
INSERT INTO accounts.spotify_users (spotify_id, display_name, email, profile_picture)
VALUES ($1, $2, $3, $4)
RETURNING *`

func (repo *accountsDatabaseRepository) SpotifyUserInsert(ctx context.Context, spotifyUser *domain.SpotifyUser) error {
	err := repo.Insert(ctx, spotifyUser, qSpotifyUserInsert,
		spotifyUser.SpotifyID, spotifyUser.DisplayName, spotifyUser.Email, spotifyUser.ProfilePicture)
	return domain.Wrap(err, "cannot insert spotify-user", nil)
}

/////

const qSpotifyUserUpdate = ` 
UPDATE accounts.spotify_users
SET %s
WHERE %s = $1`

func (repo *accountsDatabaseRepository) SpotifyUserUpdate(ctx context.Context, key string, value any, updates map[string]any) error {
	err := repo.Update(ctx, updates, fmt.Sprintf(qSpotifyUserUpdate, "%s", key), value)
	return domain.Wrap(err, "cannot update spotify-user", nil)
}

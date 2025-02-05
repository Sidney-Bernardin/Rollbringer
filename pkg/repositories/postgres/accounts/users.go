package database

import (
	"context"
	"rollbringer/pkg/domain"

	"github.com/jmoiron/sqlx"
)

const qUserInsert = ` 
WITH inserted_user AS (
	INSERT INTO accounts.users (username)
	VALUES ($1) 
	ON CONFLICT (username)
		DO UPDATE SET username = EXCLUDED.username
	RETURNING *
)
SELECT * FROM inserted_user`

func (repo *accountsDatabaseRepository) UserInsert(ctx context.Context, user *domain.User) error {
	err := sqlx.GetContext(ctx, repo.TX, user, qUserInsert,
		user.Username)
	return domain.Wrap(err, "cannot insert user", nil)
}

/////

const qGoogleUserInsert = ` 
WITH inserted_google_user AS (

	INSERT INTO accounts.google_users (user_id, google_id, given_name, profile_picture)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (google_id) DO UPDATE SET google_id = EXCLUDED.google_id
	RETURNING *
)
SELECT * FROM inserted_google_user`

func (repo *accountsDatabaseRepository) GoogleUserInsert(ctx context.Context, googleUser *domain.GoogleUser) error {
	err := sqlx.GetContext(ctx, repo.TX, googleUser, qGoogleUserInsert,
		googleUser.UserID, googleUser.GoogleID, googleUser.GivenName, googleUser.ProfilePicture)
	return domain.Wrap(err, "cannot insert google-user", nil)
}

/////

const qSpotifyUserInsert = ` 
WITH inserted_spotify_user AS (
	INSERT INTO accounts.spotify_users (user_id, spotify_id, display_name, profile_picture)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (spotify_id) DO UPDATE SET spotify_id = EXCLUDED.spotify_id
	RETURNING *
)
SELECT * FROM inserted_spotify_user`

func (repo *accountsDatabaseRepository) SpotifyUserInsert(ctx context.Context, spotifyUser *domain.SpotifyUser) error {
	err := sqlx.GetContext(ctx, repo.TX, spotifyUser, qSpotifyUserInsert,
		spotifyUser.UserID, spotifyUser.SpotifyID, spotifyUser.DisplayName, spotifyUser.ProfilePicture)
	return domain.Wrap(err, "cannot insert spotify-user", nil)
}

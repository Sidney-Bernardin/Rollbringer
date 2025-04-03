package accounts

type spotifyUser struct {
	SpotifyID string `db:"spotify_id"`

	DisplayName    string  `db:"given_name"`
	Email          string  `db:"email"`
	ProfilePicture *string `db:"profile_picture"`
}

const (
	qSpotifyUserInsert = `
		INSERT INTO accounts.spotify_users (spotify_id, display_name, email, profile_picture)
		VALUES ($1, $2, $3, $4)`

	qSpotifyUserUpdateByID = `
		UPDATE accounts.spotify_users 
		SET 
			display_name = $2,
			email = $3,
			profile_picture = $4
		WHERE spotify_id = $1`
)

package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"rollbringer/pkg/domain"
)

/////

const qSessionInsert = ` 
WITH inserted_session AS (
	INSERT INTO accounts.sessions (user_id, csrf_token)
	VALUES ($1, $2)
	ON CONFLICT (user_id) DO UPDATE SET
		id = EXCLUDED.id,
		csrf_token = EXCLUDED.csrf_token
	RETURNING *
)
SELECT * FROM inserted_session`

func (repo *accountsDatabaseRepository) SessionInsert(ctx context.Context, session *domain.Session) error {
	err := sqlx.GetContext(ctx, repo.TX, session, qSessionInsert,
		session.UserID, session.CSRFToken)
	return domain.Wrap(err, "cannot insert session", nil)
}

/////

const qSessionGet = ` 
SELECT
	sessions.id, sessions.user_id, sessions.csrf_token,

	users.id AS "user.id",
	users.google_id AS "user.google_id",
	users.spotify_id AS "user.spotify_id",
	users.username AS "user.username",
	users.profile_picture AS "user.profile_picture",

	COALESCE(google_users.google_id, '<null>') AS "google_user.google_id",
	COALESCE(google_users.email, '<null>') AS "google_user.email",

	COALESCE(spotify_users.spotify_id, '<null>') AS "spotify_user.spotify_id",
	COALESCE(spotify_users.email, '<null>') AS "spotify_user.email"
FROM accounts.sessions
LEFT JOIN accounts.users ON sessions.user_id = users.id
LEFT JOIN accounts.google_users ON users.google_id = google_users.google_id
LEFT JOIN accounts.spotify_users ON users.spotify_id = spotify_users.spotify_id
WHERE sessions.%s = $1`

func (repo *accountsDatabaseRepository) SessionGet(ctx context.Context, key string, value any) (*domain.Session, error) {
	session := &domain.Session{}
	if err := repo.Get(ctx, session, fmt.Sprintf(qSessionGet, key), value); err != nil {
		return nil, domain.Wrap(err, "cannot select session", nil)
	}

	if session.GoogleUser.GoogleID == "<null>" {
		session.GoogleUser = nil
	}

	if session.SpotifyUser.SpotifyID == "<null>" {
		session.SpotifyUser = nil
	}

	return session, nil
}

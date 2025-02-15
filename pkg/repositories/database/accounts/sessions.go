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
	users.id, users.google_id AS "user.google_id", users.spotify_id AS "user.spotify_id", users.username AS "user.username", users.profile_picture AS "user.profile_picture",
	google_users.google_id AS "google_user.google_id", google_users.email AS "google_user.email"
FROM accounts.sessions
LEFT JOIN accounts.users ON sessions.user_id = users.id
LEFT JOIN accounts.google_users ON users.google_id = google_users.google_id
WHERE sessions.%s = $1`

func (repo *accountsDatabaseRepository) SessionGet(ctx context.Context, key string, value any) (*domain.Session, error) {
	session := &domain.Session{}
	if err := repo.Get(ctx, session, fmt.Sprintf(qSessionGet, key), value); err != nil {
		return nil, domain.Wrap(err, "cannot select session", nil)
	}
	return session, nil
}

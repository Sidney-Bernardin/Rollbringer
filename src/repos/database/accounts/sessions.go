package accounts

import (
	"context"

	"rollbringer/src/domain/services/accounts"
	"rollbringer/src/repos/database"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type session struct {
	ID uuid.UUID `db:"sessions.id"`

	UserID uuid.UUID `db:"sessions.user_id"`
	user

	CSRFToken string `db:"sessions.csrf_token"`
}

func (s session) Model() *accounts.Session {
	return &accounts.Session{
		ID:        s.ID,
		UserID:    s.UserID,
		User:      s.user.Model(),
		CSRFToken: accounts.CSRFToken(s.CSRFToken),
	}
}

func (db *accountsDatabase) GetSessionBySessionID(ctx context.Context, sessionID uuid.UUID) (*accounts.Session, error) {
	_, model, err := database.Get[session](ctx, db.Pool,
		`
			SELECT 
				sessions.id AS "sessions.id",
				sessions.user_id AS "sessions.user_id",
				sessions.csrf_token AS "sessions.csrf_token",
				users.id AS "users.id",
				users.google_id AS "users.google_id",
				users.spotify_id AS "users.spotify_id",
				users.username AS "users.username",
				users.profile_picture AS "users.profile_picture"
			FROM accounts.sessions
			LEFT JOIN accounts.users ON sessions.user_id = users.id
			WHERE sessions.id = $1
		`,
		sessionID)

	return model, errors.Wrap(err, "cannot select session by ID")
}

func (db *accountsDatabase) GetSessionBySessionIDAndCSRFToken(ctx context.Context, sessionID uuid.UUID, csrfToken string) (*accounts.Session, error) {
	_, model, err := database.Get[session](ctx, db.Pool,
		`
			SELECT
				sessions.id AS "sessions.id",
				sessions.user_id AS "sessions.user_id",
				sessions.csrf_token AS "sessions.csrf_token",
				users.id AS "users.id",
				users.google_id AS "users.google_id",
				users.spotify_id AS "users.spotify_id",
				users.username AS "users.username",
				users.profile_picture AS "users.profile_picture"
			FROM accounts.sessions
			LEFT JOIN accounts.users ON sessions.user_id = users.id
			WHERE sessions.id = $1 AND sessions.csrf_token = $2
		`,
		sessionID, csrfToken)

	return model, errors.Wrap(err, "cannot select session by ID and CSRF-token")
}

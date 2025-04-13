package accounts

import (
	"context"

	"rollbringer/src"
	"rollbringer/src/repositories/database"
	"rollbringer/src/services/accounts/models"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type sessionRow struct {
	ID uuid.UUID `db:"sessions.id"`

	UserID uuid.UUID `db:"sessions.user_id"`
	userRow

	CSRFToken string `db:"sessions.csrf_token"`
}

func (r *sessionRow) Domain() *models.Session {
	if r == nil {
		return nil
	}

	return &models.Session{
		ID:        src.UUID(r.ID),
		UserID:    src.UUID(r.UserID),
		User:      r.userRow.Domain(),
		CSRFToken: models.CSRFToken(r.CSRFToken),
	}
}

func (db *accountsDatabase) GetSessionByID(ctx context.Context, sessionID src.UUID) (*models.Session, error) {
	session, err := database.Get[sessionRow](ctx, db.Tx, `
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
	`, sessionID)

	return session.Domain(), errors.Wrap(err, "cannot select session by ID")
}

func (db *accountsDatabase) GetSessionByIDAndCSRFToken(ctx context.Context, sessionID src.UUID, csrfToken models.CSRFToken) (*models.Session, error) {
	session, err := database.Get[sessionRow](ctx, db.Tx, `
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
	`, sessionID, csrfToken)

	return session.Domain(), errors.Wrap(err, "cannot select session by ID and CSRF-token")
}

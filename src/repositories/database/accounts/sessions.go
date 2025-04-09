package accounts

import (
	"context"

	"rollbringer/src"
	"rollbringer/src/services/accounts/models"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type session struct {
	ID uuid.UUID `db:"id"`

	UserID uuid.UUID `db:"user_id"`
	User   *user     `db:"user"`

	CSRFToken string `db:"csrf_token"`
}

func (s *session) domain() *models.Session {
	return &models.Session{
		SessionID: src.UUID(s.ID),
		UserID:    src.UUID(s.UserID),
		CSRFToken: models.CSRFToken(s.CSRFToken),
		User: &models.User{
			UserID:    src.UUID(s.User.ID),
			GoogleID:  s.User.GoogleID,
			SpotifyID: s.User.SpotifyID,
			Username:  models.Username(s.User.Username),
		},
	}
}

const qUspertSession = `
	INSERT INTO accounts.sessions (id, user_id, csrf_token)
	VALUES ($1, $2, $3)
	ON CONFLICT (user_id) DO UPDATE SET
		id = EXCLUDED.id,
		csrf_token = EXCLUDED.csrf_token`

/////

const (
	qSelectSessionBySessionID = `
		SELECT 
			sessions.id, sessions.user_id, sessions.csrf_token,
			users.id AS "user.id",
			users.google_id AS "user.google_id",
			users.spotify_id AS "user.spotify_id",
			users.username AS "user.username",
			users.profile_picture AS "user.profile_picture"
		FROM accounts.sessions
		LEFT JOIN accounts.users ON sessions.user_id = users.id
		WHERE sessions.id = $1`

	qSelectSessionBySessionIDAndCSRFToken = `
		SELECT
			sessions.id, sessions.user_id, sessions.csrf_token,
			users.id AS "user.id",
			users.google_id AS "user.google_id",
			users.spotify_id AS "user.spotify_id",
			users.username AS "user.username",
			users.profile_picture AS "user.profile_picture"
		FROM accounts.sessions
		LEFT JOIN accounts.users ON sessions.user_id = users.id
		WHERE sessions.id = $1 AND sessions.csrf_token = $2`
)

func (db *accountsDatabase) GetSessionByID(ctx context.Context, sessionID src.UUID) (*models.Session, error) {

	var s session
	if err := db.SelectOne(ctx, &s, qSelectSessionBySessionID, sessionID); err != nil {
		return nil, errors.Wrap(err, "cannot select session by ID")
	}

	return s.domain(), nil
}

func (db *accountsDatabase) GetSessionByIDAndCSRFToken(ctx context.Context, sessionID src.UUID, csrfToken models.CSRFToken) (*models.Session, error) {

	var s session
	if err := db.SelectOne(ctx, &s, qSelectSessionBySessionIDAndCSRFToken, sessionID, csrfToken); err != nil {
		return nil, errors.Wrap(err, "cannot select session by ID")
	}

	return s.domain(), nil
}

package database

import (
	"context"
	"database/sql"
	"rollbringer/pkg/models"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (database *Database) Login(ctx context.Context, googleID string) (*models.Session, error) {

	// Get the user with the google-ID from the database.
	var user models.User
	q := `SELECT id, google_id FROM users WHERE google_id=$1`
	err := database.db.QueryRowContext(ctx, q, googleID).Scan(&user.ID, &user.GoogleID)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.Wrap(err, "cannot select user")
	}

	if err == sql.ErrNoRows {

		// Insert a new user with the google-ID into the database.
		user.ID = uuid.New()
		q = `INSERT INTO users (id, google_id, username) VALUES ($1, $2, $3)`
		_, err := database.db.ExecContext(ctx, q, user.ID, googleID, "abc123")
		if err != nil {
			return nil, errors.Wrap(err, "cannot insert user")
		}
	}

	q = `
		INSERT INTO sessions (id, csrf_token, user_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO 
			UPDATE SET id=EXCLUDED.id, csrf_token=EXCLUDED.csrf_token, user_id=EXCLUDED.user_id`

	session := &models.Session{
		ID:        uuid.New(),
		CSRFToken: uuid.New(),
		UserID:    user.ID,
	}

	// Insert a new session for the user into the database.
	if _, err := database.db.ExecContext(ctx, q, session.ID, session.CSRFToken, user.ID); err != nil {
		return nil, errors.Wrap(err, "cannot insert session")
	}

	return session, nil
}

func (database *Database) GetUser(ctx context.Context, userID uuid.UUID) (*models.User, error) {

	var user models.User
	q := `SELECT id, username FROM users WHERE id=$1`
	err := database.db.QueryRowContext(ctx, q, userID).Scan(&user.ID, &user.Username)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}

		return nil, errors.Wrap(err, "cannot select user")
	}

	return &user, nil
}

func (database *Database) GetSession(ctx context.Context, sessionID string) (*models.Session, error) {

	var session models.Session
	q := `SELECT id, csrf_token, user_id FROM sessions WHERE id=$1`
	err := database.db.QueryRow(q, sessionID).Scan(&session.ID, &session.CSRFToken, &session.UserID)

	if err != nil && err != sql.ErrNoRows {

		if err == sql.ErrNoRows {
			return nil, ErrUnauthorized
		}

		return nil, errors.Wrap(err, "cannot select session")
	}

	return &session, nil
}

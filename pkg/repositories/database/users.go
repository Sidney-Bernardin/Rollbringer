package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
)

// InsertUser inserts the user. If the username is takan, returns
// domain.ErrUsernameTaken.
func (db *Database) InsertUser(ctx context.Context, user *domain.User) error {

	user.ID = uuid.New().String()

	// Insert the user.
	result, err := db.conn.Exec(
		`INSERT INTO users (id, username, google_id) VALUES ($1, $2, $3)
			ON CONFLICT (google_id) DO NOTHING`,
		user.ID, user.Username, user.GoogleID)

	if err != nil {
		return errors.Wrap(err, "cannot insert user")
	}

	// Get the number of rows affected.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "cannot get rows affected")
	}

	if rowsAffected == 0 {

		// Get the user-ID of the user with the google-ID.
		err := db.conn.QueryRow(
			`SELECT id FROM users WHERE google_id = $1`, user.GoogleID).
			Scan(&user.ID)

		if err != nil {
			return errors.Wrap(err, "cannot get user")
		}
	}

	return nil
}

// UpsertSession upserts the session.
func (db *Database) UpsertSession(ctx context.Context, session *domain.Session) error {

	session.ID = uuid.New().String()
	db.parseUUIDs(&session.UserID)

	// Upsert the session.
	_, err := db.conn.Exec(
		`INSERT INTO sessions (id, csrf_token, user_id) VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO 
			UPDATE SET id=EXCLUDED.id, csrf_token=EXCLUDED.csrf_token, user_id=EXCLUDED.user_id`,
		session.ID, session.CSRFToken, session.UserID)

	return errors.Wrap(err, "cannot upsert session")
}

// GetUser returns the user with the user-ID. If the user doesn't exist,
// returns domain.ErrUserNotFound.
func (db *Database) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	db.parseUUIDs(&userID)

	// Get the user with the user-ID.
	var user domain.User
	err := db.conn.QueryRow(
		`SELECT id, username, google_id FROM users WHERE id = $1`, userID).
		Scan(&user.ID, &user.Username, &user.GoogleID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}

		return nil, errors.Wrap(err, "cannot select user")
	}

	return &user, nil
}

// GetSession returns the session with the session-ID. If the session doesn't
// exist, returns domain.ErrUnauthorized.
func (db *Database) GetSession(ctx context.Context, sessionID string) (*domain.Session, error) {
	db.parseUUIDs(&sessionID)

	// Get the session with the session-ID.
	var session domain.Session
	err := db.conn.QueryRow(
		`SELECT id, user_id, csrf_token FROM sessions WHERE id = $1`, sessionID).
		Scan(&session.ID, &session.UserID, &session.CSRFToken)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUnauthorized
		}

		return nil, errors.Wrap(err, "cannot select session")
	}

	return &session, nil
}

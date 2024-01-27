package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
)

// Login inserts a new session for the user with the google-ID. If the user
// doesn't exist, a new one will be inserted.
func (db *Database) Login(ctx context.Context, googleID string) (string, error) {

	// Get a user with the google-ID.
	rows, err := db.conn.Query(ctx, `SELECT id FROM users WHERE google_id = $1`, googleID)
	if err != nil {
		return "", errors.Wrap(err, "cannot select user")
	}

	// Scan first row into a user model.
	user, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[domain.User])
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", errors.Wrap(err, "cannot scan user")
	}

	rows.Close()

	if user == nil {
		user = &domain.User{ID: uuid.New().String()}

		// Insert a new user.
		_, err = db.conn.Exec(ctx,
			`INSERT INTO users (id, username, google_id) VALUES ($1, $2, $3)`,
			user.ID, "abc123", googleID)

		if err != nil {
			return "", errors.Wrap(err, "cannot insert user")
		}
	}

	sessionID := uuid.New().String()

	// Insert a new session for the user.
	_, err = db.conn.Exec(ctx,
		`
			INSERT INTO sessions (id, csrf_token, user_id) 
				VALUES ($1, $2, $3)
			ON CONFLICT (user_id) DO 
				UPDATE SET id=EXCLUDED.id, csrf_token=EXCLUDED.csrf_token, user_id=EXCLUDED.user_id
		`,
		sessionID, uuid.New(), user.ID)

	if err != nil {
		return "", errors.Wrap(err, "cannot insert session")
	}

	return sessionID, nil
}

// GetUser returns the user with the user-ID from the database. If the user
// doesn't exist, returns domain.ErrUserNotFound.
func (db *Database) GetUser(ctx context.Context, userID string) (*domain.User, error) {

	userUUID, _ := uuid.Parse(userID)

	// Get the user with the user-ID.
	rows, err := db.conn.Query(ctx, `SELECT * FROM users WHERE id = $1`, userUUID)
	if err != nil {
		return nil, errors.Wrap(err, "cannot select user")
	}
	defer rows.Close()

	// Scan into a user model.
	user, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[domain.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}

		return nil, errors.Wrap(err, "cannot scan user")
	}

	return user, nil
}

// GetSession returns the session with the session-ID from the database. If the
// session doesn't exist, returns domain.ErrUnauthorized.
func (db *Database) GetSession(ctx context.Context, sessionID string) (*domain.Session, error) {

	sessionUUID, _ := uuid.Parse(sessionID)

	// Get the session with the session-ID.
	rows, err := db.conn.Query(ctx, `SELECT * FROM sessions WHERE id = $1`, sessionUUID)
	if err != nil {
		return nil, errors.Wrap(err, "cannot select session")
	}
	defer rows.Close()

	// Scan into a session model.
	session, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[domain.Session])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUnauthorized
		}

		return nil, errors.Wrap(err, "cannot scan session")
	}

	return session, nil
}

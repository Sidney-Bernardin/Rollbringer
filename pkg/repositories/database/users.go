package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
)

// InsertUser inserts the user. If the username is takan, returns
// domain.ErrUsernameTaken.
func (db *Database) InsertUser(ctx context.Context, user *domain.User) error {

	user.ID = uuid.New().String()

	// Insert the user.
	cmdTag, err := db.conn.Exec(ctx,
		`
			INSERT INTO users (id, username, google_id) VALUES ($1, $2, $3)
			ON CONFLICT (google_id) DO NOTHING
		`,
		user.ID, user.Username, user.GoogleID)

	if err != nil {
		return errors.Wrap(err, "cannot insert user")
	}

	if cmdTag.RowsAffected() == 0 {

		// Get the user-ID of the user with the google-ID.
		rows, err := db.conn.Query(ctx, `SELECT id FROM users WHERE google_id = $1`, user.GoogleID)
		if err != nil {
			return errors.Wrap(err, "cannot select user")
		}
		defer rows.Close()

		// Scan into the user's ID.
		rows.Next()
		if err := rows.Scan(&user.ID); err != nil {
			return errors.Wrap(err, "cannot scan user ID")
		}
	}

	return nil
}

// UpsertSession upserts the session.
func (db *Database) UpsertSession(ctx context.Context, session *domain.Session) error {

	session.ID = uuid.New().String()

	// Upsert the session.
	_, err := db.conn.Exec(ctx,
		`
			INSERT INTO sessions (id, csrf_token, user_id) VALUES ($1, $2, $3)
			ON CONFLICT (user_id) DO 
				UPDATE SET id=EXCLUDED.id, csrf_token=EXCLUDED.csrf_token, user_id=EXCLUDED.user_id
		`,
		session.ID, session.CSRFToken, session.UserID)

	return errors.Wrap(err, "cannot upsert session")
}

// GetUser returns the user with the user-ID. If the user doesn't exist,
// returns domain.ErrUserNotFound.
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

// GetSession returns the session with the session-ID. If the session doesn't
// exist, returns domain.ErrUnauthorized.
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

package database

import (
	"context"
	"database/sql"
	"rollbringer/pkg/models"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (database *Database) Login(ctx context.Context, googleID string) (*models.Session, error) {

	var user models.User

	q := `SELECT id, google_id FROM users WHERE google_id=$1`
	err := database.db.QueryRowContext(ctx, q, googleID).Scan(&user.ID, &user.GoogleID)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.Wrap(err, "cannot select user")
	}

	if err == sql.ErrNoRows {
		user = models.User{
			ID:       uuid.New(),
			GoogleID: googleID,
			Username: "abc123",
		}

		q = `INSERT INTO users (id, google_id, username) VALUES ($1, $2, $3)`
		_, err := database.db.ExecContext(ctx, q, user.ID, user.GoogleID, user.Username)
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

	if _, err := database.db.ExecContext(ctx, q, session.ID, session.CSRFToken, user.ID); err != nil {
		return nil, errors.Wrap(err, "cannot insert session")
	}

	return session, nil
}

func (database *Database) GetSession(ctx context.Context, id string) (*models.Session, error) {
	row := database.db.QueryRow(`SELECT * FROM sessions WHERE id=$1`, id)

	var session models.Session
	err := row.Scan(&session.ID, &session.CSRFToken, &session.UserID)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.Wrap(err, "cannot get session")
	}

	return &session, nil
}

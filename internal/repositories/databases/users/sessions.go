package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"rollbringer/internal"
)

var sessionViewColumns = map[internal.SessionView]string{
	internal.SessionViewAll: "sessions.id, sessions.user_id, sessions.csrf_token",
}

type dbSession struct {
	ID uuid.UUID `db:"id"`

	UserID    uuid.UUID `db:"user_id"`
	CSRFToken string    `db:"csrf_token"`
}

func (session *dbSession) internalized() *internal.Session {
	if session != nil {
		return &internal.Session{
			ID:        session.ID,
			UserID:    session.UserID,
			CSRFToken: session.CSRFToken,
		}
	}
	return nil
}

func (db *UsersDatabase) SessionUpsert(ctx context.Context, session *internal.Session) error {
	err := db.TX.QueryRowxContext(ctx,
		`INSERT INTO sessions (id, user_id, csrf_token)
			VALUES (:id, :user_id, :csrf_token)
		ON CONFLICT (user_id) DO UPDATE SET 
			id = EXCLUDED.id,
			user_id = EXCLUDED.user_id,
			csrf_token = EXCLUDED.csrf_token
		RETURNING id`,
	).Scan(&session.ID)

	return errors.Wrap(err, "cannot insert session")
}

func (db *UsersDatabase) SessionGet(ctx context.Context, sessionID uuid.UUID, view internal.SessionView) (*internal.Session, error) {
	columns, ok := sessionViewColumns[view]
	if !ok {
		return nil, fmt.Errorf("bad session view %d", view)
	}
	query := fmt.Sprintf(`SELECT %s FROM sessions WHERE id = $1`, columns)

	var session dbSession
	if err := sqlx.GetContext(ctx, db.TX, &session, query, sessionID); err != nil {
		if err == sql.ErrNoRows {
			return nil, internal.NewProblemDetail(ctx, &internal.PDOptions{
				Type:   internal.PDTypeSessionNotFound,
				Detail: "Can't find a session with the given session-ID.",
				Extra: map[string]any{
					"session_id": sessionID,
				},
			})
		}

		return nil, errors.Wrap(err, "cannot select session")
	}

	return session.internalized(), nil
}

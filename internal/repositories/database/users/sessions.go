package users

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"rollbringer/internal"
	"rollbringer/internal/repositories/database"
)

func sessionColumns(view internal.SessionView) string {
	switch view {
	case internal.SessionViewPage:
		return `sessions.*,` +
			`users.id AS "user.id",` +
			`users.username AS "user.username",` +
			`users.google_id AS "user.google_id"`

	default:
		return `sessions.*`
	}
}

func sessionJoins(view internal.SessionView) string {
	switch view {
	case internal.SessionViewPage:
		return `LEFT JOIN users.users ON users.id = sessions.user_id`
	default:
		return ``
	}
}

func (db *usersSchema) SessionUpsert(ctx context.Context, session *internal.Session) error {
	err := db.TX.QueryRowxContext(ctx,
		`INSERT INTO users.sessions (id, user_id, csrf_token)
			VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE SET 
			id = EXCLUDED.id,
			user_id = EXCLUDED.user_id,
			csrf_token = EXCLUDED.csrf_token
		RETURNING id`,
		uuid.New(), session.UserID, session.CSRFToken,
	).Scan(&session.ID)

	return errors.Wrap(err, "cannot insert session")
}

func (db *usersSchema) SessionGet(ctx context.Context, sessionID uuid.UUID, view internal.SessionView) (*internal.Session, error) {

	var session database.Session
	query := fmt.Sprintf(`SELECT %s FROM users.sessions %s WHERE sessions.id = $1`, sessionColumns(view), sessionJoins(view))
	if err := sqlx.GetContext(ctx, db.TX, &session, query, sessionID); err != nil {
		if err == sql.ErrNoRows {
			return nil, internal.NewProblemDetail(ctx, internal.PDOpts{
				Type:   internal.PDTypeSessionNotFound,
				Detail: "Cannot find a session with the given session_id.",
				Extra: map[string]any{
					"session_id": sessionID,
				},
			})
		}

		return nil, errors.Wrap(err, "cannot select session")
	}

	return session.Internalized(), nil
}

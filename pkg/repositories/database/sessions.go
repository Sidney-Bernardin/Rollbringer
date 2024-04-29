package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
)

var sessionViewColumns = map[domain.SessionView]string{
	domain.SessionViewAll: "sessions.id, sessions.user_id, sessions.csrf_token",
}

type sessionModel struct {
	ID uuid.UUID `db:"id"`

	UserID    uuid.UUID `db:"user_id"`
	CSRFToken string    `db:"csrf_token"`
}

func (session *sessionModel) domain() *domain.Session {
	if session != nil {
		return &domain.Session{
			ID:        session.ID,
			UserID:    session.UserID,
			CSRFToken: session.CSRFToken,
		}
	}
	return nil
}

func (db *Database) UpsertSession(ctx context.Context, session *domain.Session) error {

	model := sessionModel{
		ID:        uuid.New(),
		UserID:    session.UserID,
		CSRFToken: session.CSRFToken,
	}

	// Upsert the session.
	rows, err := sqlx.NamedQueryContext(ctx, db.tx,
		`
			INSERT INTO sessions (id, user_id, csrf_token)
				VALUES (:id, :user_id, :csrf_token)
			ON CONFLICT (user_id) DO UPDATE SET 
				id = EXCLUDED.id,
				user_id = EXCLUDED.user_id,
				csrf_token = EXCLUDED.csrf_token
			RETURNING id
		`,
		model,
	)

	if err != nil {
		return errors.Wrap(err, "cannot insert session")
	}
	defer rows.Close()

	rows.Next()

	// Scan the row into the model.
	if err := rows.StructScan(&model); err != nil {
		return errors.Wrap(err, "cannot scan user")
	}

	*session = *model.domain()
	return nil
}

func (db *Database) GetSession(ctx context.Context, sessionID uuid.UUID, view domain.SessionView) (*domain.Session, error) {

	// Build a query to select a session with the session-ID.
	query := fmt.Sprintf(
		`SELECT %s FROM sessions WHERE id = $1`,
		sessionViewColumns[view],
	)

	// Execute the query.
	var model sessionModel
	if err := sqlx.GetContext(ctx, db.tx, &model, query, sessionID); err != nil {
		if err == sql.ErrNoRows {
			return nil, &domain.ProblemDetail{
				Type:   domain.PDTypeSessionNotFound,
				Detail: fmt.Sprintf("Cannot find a session with the session-ID"),
			}
		}

		return nil, errors.Wrap(err, "cannot select session")
	}

	return model.domain(), nil
}

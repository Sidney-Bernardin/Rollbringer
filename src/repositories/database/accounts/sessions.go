package accounts

import (
	"context"
	"fmt"

	"rollbringer/src"
	"rollbringer/src/domain"
	"rollbringer/src/domain/accounts"
	"rollbringer/src/repositories/database"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type session struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	CSRFToken string    `db:"csrf_token"`

	User *user `db:"user"`
}

const (
	qSessionUpsert = `
		INSERT INTO accounts.sessions (id, user_id, csrf_token)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE SET
			id = EXCLUDED.id,
			csrf_token = EXCLUDED.csrf_token`

	qSessionSelectByID             = `SELECT %s FROM accounts.sessions %s WHERE sessions.id = $1`
	qSessionSelectByIDAndCSRFToken = `SELECT %s FROM accounts.sessions %s WHERE sessions.id = $1 AND sessions.csrf_token = $2`
)

func (db *accountsDatabase) querySession(ctx context.Context, crudFunc database.CRUDFunc, view any, query string, args ...any) error {

	var columns, joins string
	switch view.(type) {
	case *accounts.ViewSessionInfo:
		columns = ` sessions.id, sessions.user_id, sessions.csrf_token,
			users.id AS "user.id",
			users.google_id AS "user.google_id",
			users.spotify_id AS "user.spotify_id",
			users.username AS "user.username",
			users.profile_picture AS "user.profile_picture"`
		joins = `LEFT JOIN accounts.users ON sessions.user_id = users.id`
	default:
		return &src.ExternalError{Type: domain.ExternalErrorTypeViewInvalid}
	}

	var s session
	if err := crudFunc(ctx, &s, fmt.Sprintf(query, columns, joins), args...); err != nil {
		return err
	}

	switch v := view.(type) {
	case *accounts.ViewSessionInfo:
		*v = accounts.ViewSessionInfo{
			SessionID: s.ID.String(),
			UserID:    s.UserID.String(),
			CSRFToken: s.CSRFToken,
			UserInfo: accounts.ViewUserInfo{
				UserID:    s.User.ID.String(),
				GoogleID:  s.User.GoogleID,
				SpotifyID: s.User.SpotifyID,
				Username:  s.User.Username,
			},
		}
	}

	return nil
}

func (db *accountsDatabase) SessionGetByID(ctx context.Context, view any, sessionID domain.UUID) error {
	id, _ := uuid.Parse(sessionID.String())
	err := db.querySession(ctx, db.CRUDGet, view, qSessionSelectByID, id)
	return errors.Wrap(err, "cannot get session by ID")
}

func (db *accountsDatabase) SessionGetByIDAndCSRFToken(ctx context.Context, view any, sessionID domain.UUID, csrfToken accounts.CSRFToken) error {
	id, _ := uuid.Parse(sessionID.String())
	err := db.querySession(ctx, db.CRUDGet, view, qSessionSelectByIDAndCSRFToken, id, csrfToken)
	return errors.Wrap(err, "cannot get session by ID and CSRF-token")
}

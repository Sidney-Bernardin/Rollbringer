package accounts

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"rollbringer/src/repositories/database"
)

type session struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	CSRFToken string    `db:"csrf_token"`
}

const (
	qSessionInsert = `
		INSERT INTO accounts.google_users (google_id, given_name, email, profile_picture)
		VALUES ($1, $2, $3, $4)`
)

func (db *accountsDatabase) querySession(ctx context.Context, crudFunc database.CRUDFunc, view any, query string, args ...any) error {

	var columns, joins string
	switch view.(type) {
	default:
		columns = `sessions.id, sessions.user_id, sessions.csrf_token`
	}

	var s session
	if err := crudFunc(ctx, &s, fmt.Sprintf(query, columns, joins), args...); err != nil {
		return err
	}

	switch v := view.(type) {
	case *uuid.UUID:
		*v = s.ID
	}

	return nil
}

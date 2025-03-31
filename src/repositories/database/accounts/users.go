package accounts

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"rollbringer/src/domain/accounts"
	"rollbringer/src/repositories/database"
)

type user struct {
	ID       uuid.UUID `db:"id"`
	Username string    `db:"username"`
}

const (
	qInsertUser = `
		WITH inserted_user AS (
			INSERT INTO accounts.users (id, username)
			VALUES ($1, $2)
			RETURNING *
		)
		SELECT %s FROM inserted_user %s`

	qUserSelectByGoogleID = `
		SELECT %s FROM accounts.users %s
		WHERE google_id = $1`

	qUserSelectByUsername = `
		SELECT %s FROM accounts.users %s
		WHERE username = $1`
)

func (db *accountsDatabase) queryUser(ctx context.Context, crudFunc database.CRUDFunc, view any, query string, args ...any) error {

	var columns, joins string
	switch view.(type) {
	default:
		columns = `users.id, users.username`
	}

	var u user
	if err := crudFunc(ctx, &u, fmt.Sprintf(query, columns, joins), args...); err != nil {
		return err
	}

	switch v := view.(type) {
	case *uuid.UUID:
		*v = u.ID
	case *accounts.ViewUserInfo:
		v.UserID = u.ID.String()
		v.Username = u.Username
	}

	return nil
}

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

func (db *accountsDatabase) userQuery(ctx context.Context, crudFunc database.CRUDFunc, view any, query string, args ...any) error {

	var columns, joins string
	switch view.(type) {
	default:
		columns = `users.id, users.username`
	}

	var r user
	if err := crudFunc(ctx, &r, fmt.Sprintf(query, columns, joins), args...); err != nil {
		return err
	}

	switch v := view.(type) {
	case *accounts.UserInfo:
		v.UserID = r.ID.String()
		v.Username = r.Username
	}

	return nil
}

/////

const qUserInsert = `
	WITH inserted_user AS (
		INSERT INTO accounts.users (username)
		VALUES ($1)
		RETURNING *
	)
	SELECT %s FROM inserted_user %s`

func (db *accountsDatabase) UserCreate(ctx context.Context, view any, cmd *accounts.CmdUserCreate) error {
	return db.userQuery(ctx, db.CRUDInsert, view, qUserInsert, cmd.Username)
}

/////

const qUserGetByUsername = `SELECT %s FROM accounts.users %s WHERE username = $1`

func (db *accountsDatabase) UserGetByUsername(ctx context.Context, view any, username accounts.Username) error {
	return db.userQuery(ctx, db.CRUDGet, view, qUserGetByUsername, username)
}

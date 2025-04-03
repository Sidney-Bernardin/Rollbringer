package accounts

import (
	"context"
	"fmt"

	"rollbringer/src"
	"rollbringer/src/domain"
	"rollbringer/src/domain/accounts"
	"rollbringer/src/repositories/database"

	"github.com/google/uuid"
)

type user struct {
	ID             uuid.UUID `db:"id"`
	GoogleID       string    `db:"google_id"`
	SpotifyID      string    `db:"spotify_id"`
	Username       string    `db:"username"`
	ProfilePicture string    `db:"profile_picture"`
}

const (
	qInsertUser = `
		INSERT INTO accounts.users (id, google_id, spotify_id, username, profile_picture)
		VALUES ($1, $2, $3, $4, $5)`

	qUserSelectByGoogleID = `
		SELECT %s FROM accounts.users %s
		WHERE google_id = $1`

	qUserSelectBySpotifyID = `
		SELECT %s FROM accounts.users %s
		WHERE spotify_id = $1`

	qUserSelectByUsername = `
		SELECT %s FROM accounts.users %s
		WHERE username = $1`
)

func (db *accountsDatabase) queryUser(ctx context.Context, crudFunc database.CRUDFunc, view any, query string, args ...any) error {

	var columns, joins string
	switch view.(type) {
	case *domain.UUID:
		columns = `users.id`
	case *accounts.ViewUserInfo:
		columns = `users.id, users.username, users.profile_picture`
	default:
		return &src.ExternalError{Type: domain.ExternalErrorTypeViewInvalid}
	}

	var u user
	if err := crudFunc(ctx, &u, fmt.Sprintf(query, columns, joins), args...); err != nil {
		return err
	}

	switch v := view.(type) {
	case *domain.UUID:
		*v = domain.UUID(u.ID)
	case *accounts.ViewUserInfo:
		v.UserID = u.ID.String()
		v.Username = u.Username
		v.ProfilePicture = u.ProfilePicture
	}

	return nil
}

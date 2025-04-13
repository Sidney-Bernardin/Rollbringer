package accounts

import (
	"context"
	"rollbringer/src"
	"rollbringer/src/repositories/database"
	"rollbringer/src/services/accounts/models"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
)

type userRow struct {
	ID             pgtype.UUID `db:"users.id"`
	GoogleID       *string     `db:"users.google_id"`
	SpotifyID      *string     `db:"users.spotify_id"`
	Username       string      `db:"users.username"`
	ProfilePicture string      `db:"users.profile_picture"`
}

func (r *userRow) Domain() *models.User {
	if r == nil {
		return nil
	}

	return &models.User{
		ID:             src.UUID(r.ID.Bytes),
		GoogleID:       r.GoogleID,
		SpotifyID:      r.SpotifyID,
		Username:       models.Username(r.Username),
		ProfilePicture: r.ProfilePicture,
	}
}

type usersByRoomRow struct {
	RoomID          pgtype.UUID   `db:"room_users.room_id"`
	UserIDs         []pgtype.UUID `db:"users.user_ids"`
	GoogleIDs       []*string     `db:"users.google_ids"`
	SpotifyIDs      []*string     `db:"users.spotify_ids"`
	Usernames       []string      `db:"users.usernames"`
	ProfilePictures []string      `db:"users.profile_pictures"`
}

func (r *usersByRoomRow) Domain() []*models.User {
	if r == nil {
		return nil
	}

	users := make([]*models.User, len(r.UserIDs))
	for i := range len(users) {
		users[i] = &models.User{
			ID:             src.UUID(r.UserIDs[i].Bytes),
			GoogleID:       r.GoogleIDs[i],
			SpotifyID:      r.SpotifyIDs[i],
			Username:       models.Username(r.Usernames[i]),
			ProfilePicture: r.ProfilePictures[i],
		}
	}

	return users
}

func (db *accountsDatabase) GetUsersByRoomIDs(ctx context.Context, roomIDs ...src.UUID) (map[src.UUID][]*models.User, error) {
	var res = map[src.UUID][]*models.User{}
	if len(roomIDs) < 1 {
		return res, nil
	}

	rows, err := database.Gets[usersByRoomRow](ctx, db.Tx, `
		SELECT
			room_users.room_id AS "room_users.room_id",
			json_agg(users.id) AS "users.user_ids",
			json_agg(users.google_id) AS "users.google_ids",
			json_agg(users.spotify_id) AS "users.spotify_ids",
			json_agg(users.username) AS "users.usernames",
			json_agg(users.profile_picture) AS "users.profile_pictures"
		FROM accounts.users
		LEFT JOIN room_users ON users.id = room_users.user_id
		WHERE room_users.room_id = ANY($1)
		GROUP BY room_users.room_id
	`, roomIDs)
	if err != nil {
		return nil, errors.Wrap(err, "cannot select users by room-IDs")
	}

	for _, row := range rows {
		res[src.UUID(row.RoomID.Bytes)] = row.Domain()
	}
	return res, nil
}

package database

import (
	"context"
	"rollbringer/pkg/domain"

	"github.com/jmoiron/sqlx"
)

/////

const qSessionInsert = ` 
WITH inserted_session AS (
	INSERT INTO accounts.sessions (user_id, csrf_token)
	VALUES ($1, $2)
	ON CONFLICT (user_id) DO UPDATE SET
		id = EXCLUDED.id,
		csrf_token = EXCLUDED.csrf_token
	RETURNING *
)
SELECT * FROM inserted_session`

func (repo *accountsDatabaseRepository) SessionInsert(ctx context.Context, session *domain.Session) error {
	err := sqlx.GetContext(ctx, repo.TX, session, qSessionInsert,
		session.UserID, session.CSRFToken)
	return domain.Wrap(err, "cannot insert session", nil)
}

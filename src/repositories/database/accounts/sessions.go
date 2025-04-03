package accounts

import "rollbringer/src/domain"

type session struct {
	ID        domain.UUID `db:"id"`
	UserID    domain.UUID `db:"user_id"`
	CSRFToken string      `db:"csrf_token"`
}

const (
	qSessionUpsert = `
		INSERT INTO accounts.sessions (id, user_id, csrf_token)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE SET
			id = EXCLUDED.id,
			csrf_token = EXCLUDED.csrf_token`
)

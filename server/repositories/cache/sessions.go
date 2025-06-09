package cache

import (
	"context"
	"fmt"

	"github.com/Sidney-Bernardin/Rollbringer/server"

	"github.com/pkg/errors"
)

type Session struct {
	ID        server.UUID `redis:"id"`
	UserID    server.UUID `redis:"user_id"`
	CSRFToken string      `redis:"csrf_token"`
}

func (cache *Cache) SetSession(ctx context.Context, userID server.UUID) (server.UUID, error) {

	session := &Session{
		ID:        server.NewRandomUUID(),
		UserID:    userID,
		CSRFToken: server.CreateRandomString(),
	}

	key := fmt.Sprintf("session:%s", session.ID)
	tx := cache.client.TxPipeline()

	if err := tx.JSONSet(ctx, key, "$", session).Err(); err != nil {
		return server.UUID{}, errors.Wrap(err, "cannot put session")
	}

	if err := tx.Expire(ctx, key, cache.config.SessionTimeout).Err(); err != nil {
		return server.UUID{}, errors.Wrap(err, "cannot expire session")
	}

	if _, err := tx.Exec(ctx); err != nil {
		return server.UUID{}, errors.Wrap(err, "cannot execute transaction")
	}

	return session.ID, nil
}

func (cache *Cache) GetSession(ctx context.Context, sessionID server.UUID) (*Session, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	session, err := jsonGet[*Session](ctx, cache.client, key, ".")
	return session, errors.WithStack(err)
}

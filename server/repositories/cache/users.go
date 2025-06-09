package cache

import (
	"context"
	"fmt"

	"github.com/Sidney-Bernardin/Rollbringer/server"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/sql/queries"

	"github.com/pkg/errors"
)

func (cache *Cache) SetUser(ctx context.Context, user *queries.SelectUserRow) error {
	key := fmt.Sprintf("user:%s", user.ID)
	err := cache.client.JSONSet(ctx, key, "$", user).Err()
	return errors.WithStack(err)
}

func (cache *Cache) GetUser(ctx context.Context, userID server.UUID) (*queries.SelectUserRow, error) {
	key := fmt.Sprintf("user:%s", userID)
	user, err := jsonGet[*queries.SelectUserRow](ctx, cache.client, key, ".")
	return user, errors.WithStack(err)
}

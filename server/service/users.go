package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/Sidney-Bernardin/Rollbringer/server"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/cache"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/sql"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/sql/queries"

	"github.com/pkg/errors"
)

func (svc *Service) GetUser(ctx context.Context, userID server.UUID) (*queries.SelectUserRow, error) {

	user, err := svc.Cache.GetUser(ctx, userID)
	if !errors.Is(err, cache.ErrNotFound) {
		return user, errors.Wrap(err, "cannot get user from Cache")
	}

	user, err = svc.SQL.SelectUser(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &server.UserError{
				Type:    server.UserErrorTypeUserNotFound,
				Message: "Cannot find a user with that ID",
			}
		}

		return nil, errors.Wrap(err, "cannot get user from SQL")
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := svc.Cache.SetUser(ctx, user); err != nil {
			svc.Log.Log(ctx, slog.LevelWarn, "Cache cannot set user",
				"err", err.Error(),
				"user_id", user.ID)
		}
	}()

	return user, nil
}

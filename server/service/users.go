package service

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/Sidney-Bernardin/Rollbringer/server"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/sql/queries"

	"github.com/pkg/errors"
)

func (svc *Service) GetUser(ctx context.Context, userID server.UUID) (*queries.GetUserRow, error) {

	user, err := svc.Nats.GetUser(ctx, userID)
	if err != nil || user != nil {
		return user, errors.Wrap(err, "cannot get user from Nats")
	}

	user, err = svc.SQL.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, server.NewUserError(server.UserErrorTypeUserNotFound, "", nil)
		}

		return nil, errors.Wrap(err, "cannot get user from SQL")
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := svc.Nats.PutUser(ctx, user); err != nil {
			svc.Log.Log(ctx, slog.LevelWarn, "Nats put user", "err", err.Error())
		}
	}()

	return user, nil
}

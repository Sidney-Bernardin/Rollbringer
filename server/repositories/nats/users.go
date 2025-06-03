package nats

import (
	"context"
	"encoding/json"

	"github.com/Sidney-Bernardin/Rollbringer/server"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/sql/queries"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/pkg/errors"
)

func (nats *Nats) PutUser(ctx context.Context, user *queries.User) error {

	userJSON, err := json.Marshal(user)
	if err != nil {
		return errors.Wrap(err, "cannot encode user")
	}

	_, err = nats.sessionsKV.Put(ctx, user.ID.String(), userJSON)
	return errors.Wrap(err, "cannot put session")
}

func (nats *Nats) GetUser(ctx context.Context, userID server.UUID) (*queries.User, error) {

	res, err := nats.sessionsKV.Get(ctx, userID.String())
	if err != nil {
		switch {
		case errors.Is(err, jetstream.ErrKeyNotFound):
			return nil, nil
		default:
			return nil, errors.Wrap(err, "cannot get user")
		}
	}

	var user queries.User
	if err := json.Unmarshal(res.Value(), &user); err != nil {
		return nil, errors.Wrap(err, "cannot decode user")
	}

	return &user, nil
}

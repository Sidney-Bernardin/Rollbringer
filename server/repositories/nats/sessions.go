package nats

import (
	"context"
	"encoding/json"

	"github.com/Sidney-Bernardin/Rollbringer/server"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/pkg/errors"
)

type Session struct {
	ID        server.UUID `json:"id"`
	UserID    server.UUID `json:"user_id"`
	CSRFToken string      `json:"csrf_token"`
}

func (nats *Nats) PutSession(ctx context.Context, userID server.UUID) (sessionID server.UUID, err error) {
	session := &Session{
		ID:        server.NewRandomUUID(),
		UserID:    userID,
		CSRFToken: server.CreateRandomString(),
	}

	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return sessionID, errors.Wrap(err, "cannot encode session")
	}

	if _, err := nats.sessionsKV.Put(ctx, session.ID.String(), sessionJSON); err != nil {
		return sessionID, errors.Wrap(err, "cannot put session")
	}

	return session.ID, nil
}

func (nats *Nats) GetSession(ctx context.Context, sessionID server.UUID) (*Session, error) {

	res, err := nats.sessionsKV.Get(ctx, sessionID.String())
	if err != nil {
		switch {
		case errors.Is(err, jetstream.ErrKeyNotFound):
			return nil, nil
		default:
			return nil, errors.Wrap(err, "cannot get session")
		}
	}

	var session Session
	if err := json.Unmarshal(res.Value(), &session); err != nil {
		return nil, errors.Wrap(err, "cannot decode session")
	}

	return &session, nil
}

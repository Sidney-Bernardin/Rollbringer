package service

import (
	"log/slog"

	"github.com/Sidney-Bernardin/Rollbringer/server"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/cache"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/cql"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/google"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/pubsub"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/sql"
)

type Service struct {
	Config *server.Config
	Log    *slog.Logger

	SQL    *sql.SQL
	CQL    *cql.CQL
	Cache  *cache.Cache
	PubSub *pubsub.PubSub
	Google *google.Google
}

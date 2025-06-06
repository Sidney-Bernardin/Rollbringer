package service

import (
	"log/slog"

	"github.com/Sidney-Bernardin/Rollbringer/server"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/google"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/nats"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/sql"
)

type Service struct {
	Config *server.Config
	Log    *slog.Logger

	SQL    *sql.SQL
	Nats   *nats.Nats
	Google *google.Google
}

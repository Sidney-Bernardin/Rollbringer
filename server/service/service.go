package service

import (
	"log/slog"

	"github.com/Sidney-Bernardin/Rollbringer/server"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/nats"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/sql"
)

type Service struct {
	config *server.Config
	log    *slog.Logger

	sql  *sql.SQL
	nats *nats.Nats
}

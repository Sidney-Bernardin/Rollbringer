package http

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/Sidney-Bernardin/Rollbringer/server"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/cassandra"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/nats"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/sql"
	"github.com/Sidney-Bernardin/Rollbringer/server/service"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"
)

var errorCodes = map[server.UserErrorType]int{
	server.UserErrorTypeInternalServerError: http.StatusInternalServerError,
}

type Server struct {
	*http.Server

	config *server.Config
	log    *slog.Logger

	svc *service.Service

	sql       *sql.SQL
	cassandra *cassandra.Cassandra
	nats      *nats.Nats
}

func (svr *Server) routes() {
	var r = chi.NewRouter()

	r.HandleFunc("/", svr.handleHome)

	svr.Server.Handler = r
}

func (svr *Server) respond(w io.Writer, r *http.Request, code int, data any) {
}

func (svr *Server) err(w io.Writer, r *http.Request, err error) {
	if err == nil {
		return
	}

	var userErr *server.UserError
	if !errors.As(err, &userErr) {
		svr.log.Error("Internal server error", "err", err.Error())
		userErr = &server.UserError{Type: server.UserErrorTypeInternalServerError}
	}

	switch w.(type) {
	case http.ResponseWriter:
		svr.respond(w, r, errorCodes[userErr.Type], userErr)
	case *websocket.Conn:
	}
}

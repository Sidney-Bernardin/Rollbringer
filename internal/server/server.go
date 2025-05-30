package server

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/Sidney-Bernardin/Rollbringer/internal"
	"github.com/Sidney-Bernardin/Rollbringer/internal/domain"
	"github.com/Sidney-Bernardin/Rollbringer/internal/repos/sql"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"
)

var errorCodes = map[domain.UserErrorType]int{
	domain.DomainErrorTypeInternalServerError: http.StatusInternalServerError,
}

type Server struct {
	*http.Server

	log    *slog.Logger
	config *internal.Config

	sql       sql.SQL
	cassandra cassandra.Cassandra
	nats      nats.Nats
}

func New(log *slog.Logger, config *internal.Config, svc domain.Service, nats nats.Nats, sql sql.SQL, cassandra cassandra.Cassandra) *Server {
	svr := &Server{
		&http.Server{
			Addr: config.HttpAddr,
		},
		log, config,
	}

	svr.routes()
	return svr
}

func (svr *Server) routes() {
	var r = chi.NewRouter()

	r.Handle("/", svr.handleHome())

	svr.Server.Handler = r
}

func (svr *Server) respond(w io.Writer, r *http.Request, code int, data any) {
}

func (svr *Server) err(w io.Writer, r *http.Request, err error) {
	if err == nil {
		return
	}

	var userErr *domain.UserError
	if !errors.As(err, &userErr) {
		svr.log.Error("Internal server error", "err", err.Error())
		userErr = &domain.UserError{Type: domain.DomainErrorTypeInternalServerError}
	}

	switch w.(type) {
	case http.ResponseWriter:
		svr.respond(w, r, errorCodes[userErr.Type], userErr)
	case *websocket.Conn:
	}
}

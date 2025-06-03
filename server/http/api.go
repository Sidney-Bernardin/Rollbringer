package http

import (
	"encoding/json"
	"io"
	"log/slog"
	"net"
	"net/http"

	"github.com/Sidney-Bernardin/Rollbringer/server"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/google"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/nats"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/sql"
	"github.com/Sidney-Bernardin/Rollbringer/server/service"

	"github.com/a-h/templ"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"
)

var errorCodes = map[server.UserErrorType]int{
	server.UserErrorTypeInternalServerError: http.StatusInternalServerError,
	server.UserErrorTypeUnauthorized:        http.StatusUnauthorized,
	server.UserErrorTypeUUIDInvalid:         http.StatusBadRequest,

	server.UserErrorTypeGoogleUserAlreadyExists: http.StatusConflict,
	server.UserErrorTypeGoogleUserNotExists:     http.StatusNotFound,
	server.UserErrorTypeUserNotFound:            http.StatusNotFound,
}

type API struct {
	*http.Server

	Config *server.Config
	Log    *slog.Logger

	Service *service.Service

	SQL    *sql.SQL
	Nats   *nats.Nats
	Google *google.Google
}

func (api *API) respond(w io.Writer, r *http.Request, code int, data any) {
	if rw, ok := w.(http.ResponseWriter); ok {
		rw.WriteHeader(code)
	}

	var err error
	switch res := data.(type) {
	case []byte:
		_, err = w.Write(res)

	case templ.Component:
		err = res.Render(r.Context(), w)
		err = errors.Wrap(err, "cannot render Templ response")

	default:
		err = json.NewEncoder(w).Encode(res)
		err = errors.Wrap(err, "cannot JSON encode response")
	}

	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, net.ErrClosed) {
		api.Log.Log(r.Context(), slog.LevelError, "Cannot write response", "err", err.Error())
	}
}

func (api *API) err(w io.Writer, r *http.Request, err error) {
	if err == nil {
		return
	}

	var userErr *server.UserError
	if !errors.As(err, &userErr) {
		api.Log.Error("Internal server error", "err", err.Error())
		userErr = &server.UserError{Type: server.UserErrorTypeInternalServerError}
	}

	switch w.(type) {
	case http.ResponseWriter:
		api.respond(w, r, errorCodes[userErr.Type], userErr)
	case *websocket.Conn:
	}
}

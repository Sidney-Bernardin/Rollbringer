package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"

	"rollbringer/src"
	"rollbringer/src/api/views"
	"rollbringer/src/domain/accounts"
	"rollbringer/src/domain/play"
)

var errCodes = map[src.ExternalErrorType]int{
	src.ExternalErrorTypeUUIDInvalid:          http.StatusBadRequest,
	src.ExternalErrorTypeEntityNotFound:       http.StatusBadRequest,
	play.ExternalErrorTypeRoomNameInvalid:     http.StatusBadRequest,
	play.ExternalErrorTypeRoomNameTaken:       http.StatusConflict,
	accounts.ExternalErrorTypeUsernameInvalid: http.StatusBadRequest,
	accounts.ExternalErrorTypeUsernameTaken:   http.StatusConflict,
}

func (svr *server) respond(w io.Writer, r *http.Request, statusCode int, res any) {

	var (
		state = svr.state(r)
		ctx   = r.Context()
	)

	if rw, ok := w.(http.ResponseWriter); ok {
		rw.WriteHeader(statusCode)
		state["status_code"] = statusCode
	}

	var err error
	switch res := res.(type) {

	case templ.Component:
		err = res.Render(r.Context(), w)
		err = errors.Wrap(err, "cannot render Templ response")

	default:
		err = json.NewEncoder(w).Encode(res)
		err = errors.Wrap(err, "cannot JSON encode response")
	}

	if err != nil {
		svr.logServerError(ctx, err)
	}
}

func (svr *server) err(w io.Writer, r *http.Request, err error) {

	var (
		state = svr.state(r)
		ctx   = r.Context()
	)

	var extErr *src.ExternalError
	if !errors.As(err, &extErr) {
		svr.logServerError(ctx, err)
		svr.respond(w, r, http.StatusInternalServerError, &views.ProblemDetail{
			Instance: state["instance"].(string),
			Type:     "internal_server_error",
		})
		return
	}

	problemDetail := &views.ProblemDetail{
		Instance: state["instance"].(string),
		Type:     string(domainErr.Type),
		Detail:   domainErr.Description,
	}

	switch w.(type) {
	case *websocket.Conn:
		svr.respond(w, r, 0, views.WebSocketMessage{
			Type:    views.WSMsgTypeError,
			Payload: problemDetail,
		})

	case http.ResponseWriter:
		svr.respond(w, r, errCodes[domainErr.Type], problemDetail)
	}
}

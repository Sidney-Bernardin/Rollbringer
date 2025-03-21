package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"

	"rollbringer/src/api/views"
	"rollbringer/src/domain"
	"rollbringer/src/domain/play"
)

var errCodes = map[domain.DomainErrorType]int{
	domain.DomainErrorTypeUUIDInvalid:   http.StatusBadRequest,
	play.DomainErrorTypeRoomNameInvalid: http.StatusBadRequest,
	play.DomainErrorTypeRoomNameTaken:   http.StatusConflict,
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

	var domainErr *domain.DomainError
	if !errors.As(err, &domainErr) {
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

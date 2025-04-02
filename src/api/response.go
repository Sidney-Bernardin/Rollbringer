package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"

	"rollbringer/src"
	"rollbringer/src/api/views"
	"rollbringer/src/domain"
	"rollbringer/src/domain/accounts"
	"rollbringer/src/domain/play"
)

func (svr *server) respond(w io.Writer, r *http.Request, statusCode int, res any) {
	var ctx = r.Context()

	if rw, ok := w.(http.ResponseWriter); ok {
		rw.WriteHeader(statusCode)

		ctx = context.WithValue(ctx, "status_code", statusCode)
		r = r.WithContext(ctx)
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

var errCodes = map[src.ExternalErrorType]int{
	externalErrorTypeInternalError:   http.StatusInternalServerError,
	externalErrorTypeUnauthorized:    http.StatusUnauthorized,
	externalErrorTypeInvalidProvider: http.StatusBadRequest,

	domain.ExternalErrorTypeUUIDInvalid: http.StatusBadRequest,
	domain.ExternalErrorTypeViewInvalid: http.StatusBadRequest,

	accounts.ExternalErrorTypeUnauthorized:          http.StatusUnauthorized,
	accounts.ExternalErrorTypeUserWithoutProviders:  http.StatusBadRequest,
	accounts.ExternalErrorTypeUsernameInvalid:       http.StatusBadRequest,
	accounts.ExternalErrorTypeUsernameTaken:         http.StatusConflict,
	accounts.ExternalErrorTypeProviderNotLinked:     http.StatusBadRequest,
	accounts.ExternalErrorTypeProviderAlreadyLinked: http.StatusBadRequest,

	play.ExternalErrorTypeRoomNotFound:    http.StatusNotFound,
	play.ExternalErrorTypeRoomNameInvalid: http.StatusBadRequest,
}

func (svr *server) err(w io.Writer, r *http.Request, err error) {
	var ctx = r.Context()

	var externalErr *src.ExternalError
	if !errors.As(err, &externalErr) {
		svr.logServerError(ctx, err)
		svr.respond(w, r, http.StatusInternalServerError, &src.ExternalError{
			Type: externalErrorTypeInternalError,
		})
		return
	}

	switch w.(type) {
	case *websocket.Conn:
		svr.respond(w, r, 0, views.WebSocketMessage{
			Type:    views.WSMsgTypeError,
			Payload: externalErr,
		})

	case http.ResponseWriter:
		svr.respond(w, r, errCodes[externalErr.Type], externalErr)
	}
}

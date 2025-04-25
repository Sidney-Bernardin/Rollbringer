package api

import (
	"encoding/json"
	"io"
	"net"
	"net/http"

	"github.com/a-h/templ"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"

	"rollbringer/src"
	"rollbringer/src/api/views"
	"rollbringer/src/services/accounts"
	account_models "rollbringer/src/services/accounts/models"
	play_models "rollbringer/src/services/play/models"
)

func (svr *server) respond(w io.Writer, r *http.Request, statusCode int, res any) {
	if rw, ok := w.(http.ResponseWriter); ok {
		rw.WriteHeader(statusCode)
	}

	var err error
	switch res := res.(type) {
	case []byte:
		_, err = w.Write(res)

	case templ.Component:
		err = res.Render(r.Context(), w)
		err = errors.Wrap(err, "cannot render Templ response")

	default:
		err = json.NewEncoder(w).Encode(res)
		err = errors.Wrap(err, "cannot marshal response")
	}

	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, net.ErrClosed) {
		svr.logServerError(r.Context(), err)
	}
}

var errCodes = map[src.ExternalErrorType]int{
	externalErrorTypeInvalidProvider: http.StatusBadRequest,

	src.ExternalErrorTypeInternalError: http.StatusInternalServerError,
	src.ExternalErrorTypeUnauthorized:  http.StatusUnauthorized,
	src.ExternalErrorTypeInvalidUUID:   http.StatusUnprocessableEntity,

	account_models.ExternalErrorTypeInvalidUsername: http.StatusBadRequest,
	accounts.ExternalErrorTypeProviderNotLinked:     http.StatusBadRequest,
	accounts.ExternalErrorTypeProviderAlreadyLinked: http.StatusConflict,

	play_models.ExternalErrorTypeInvalidRoomName: http.StatusBadRequest,
}

func (svr *server) err(w io.Writer, r *http.Request, err error) {
	if err == nil {
		return
	}

	var ctx = r.Context()

	var externalErr *src.ExternalError
	if !errors.As(err, &externalErr) {
		svr.logServerError(ctx, err)
		externalErr = &src.ExternalError{Type: src.ExternalErrorTypeInternalError}
	}

	switch w.(type) {
	case *websocket.Conn:
		svr.respond(w, r, 0, &views.WebSocketResponse{
			Operation: "error",
			Payload:   externalErr,
		})

	case http.ResponseWriter:
		svr.respond(w, r, errCodes[externalErr.Type], externalErr)
	}
}

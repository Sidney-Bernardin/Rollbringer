package api

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"rollbringer/src"
	"rollbringer/src/api/views"
	"rollbringer/src/domain"
	"rollbringer/src/domain/services/accounts"
	"rollbringer/src/domain/services/play"

	"github.com/a-h/templ"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"
)

func (svr *server) respond(w io.Writer, r *http.Request, statusCode int, res any) {
	if rw, ok := w.(http.ResponseWriter); ok {
		rw.WriteHeader(statusCode)
	}

	// Response as bytes, JSON, or a Templ component.
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
		svr.log.Log(r.Context(), src.LevelError, "Cannot write response", "err", err.Error())
	}
}

var errCodes = map[domain.ExternalErrorType]int{
	externalErrorTypeInvalidProvider: http.StatusBadRequest,

	domain.ExternalErrorTypeInternalError: http.StatusInternalServerError,
	domain.ExternalErrorTypeUnauthorized:  http.StatusUnauthorized,
	domain.ExternalErrorTypeInvalidUUID:   http.StatusUnprocessableEntity,

	accounts.ExternalErrorTypeInvalidUsername:       http.StatusBadRequest,
	accounts.ExternalErrorTypeProviderNotLinked:     http.StatusBadRequest,
	accounts.ExternalErrorTypeProviderAlreadyLinked: http.StatusConflict,

	play.ExternalErrorTypeInvalidRoomName: http.StatusBadRequest,
}

func (svr *server) err(w io.Writer, r *http.Request, err error) {
	if err == nil {
		return
	}

	// If the error isn't a domain.ExternalError, treat it as a server error.
	var externalErr *domain.ExternalError
	if !errors.As(err, &externalErr) {
		svr.log.Log(r.Context(), src.LevelError, "Internal server error", "err", err.Error())
		externalErr = &domain.ExternalError{Type: domain.ExternalErrorTypeInternalError}
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

func (svr *server) nextMessage(conn *websocket.Conn, payloadCallback func(string) any) (any, error) {

	// Receive the next message.
	var msg []byte
	if err := websocket.Message.Receive(conn, &msg); err != nil {
		return nil, errors.Wrap(err, "cannot read from WebSocket connection")
	}

	// Decode the message's operation.
	var head struct {
		Operation string `json:"operation"`
	}
	if err := json.Unmarshal(msg, &head); err != nil {
		return nil, nil
	}

	// Decode the message's payload based on it's operation.
	payload := payloadCallback(head.Operation)
	if err := json.Unmarshal(msg, payload); err != nil {
		return nil, nil
	}

	return payload, nil
}

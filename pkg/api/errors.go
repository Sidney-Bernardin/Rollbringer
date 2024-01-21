package api

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/net/websocket"

	"rollbringer/pkg/repositories/database"
)

var (
	errUnauthorized = errors.New("unauthorized")
)

func (api *API) err(writer io.Writer, e error, httpStatus, wsStatus int) {
	switch w := writer.(type) {
	case http.ResponseWriter:

		if httpStatus >= 500 {
			api.Logger.Error().Stack().Err(e).Msg("Internal server error")
			e = errors.New("internal server error")
		}

		http.Error(w, e.Error(), httpStatus)

	case *websocket.Conn:
		if wsStatus == 1011 {
			api.Logger.Error().Stack().Err(e).Msg("Internal server error")
			e = errors.New("internal server error")
		}

		if err := w.WriteClose(wsStatus); err != nil {
			api.Logger.Error().Stack().Err(err).Msg("Cannot write close status")
		}
	}
}

func (api *API) dbErr(writer io.Writer, err error) {

	var (
		res = errors.Cause(err)

		httpStatus int
		wsStatus   int
	)

	switch res {
	case database.ErrUnauthorized:
		httpStatus = http.StatusUnauthorized
		wsStatus = wsStatusPolicyViolation

	case database.ErrUserNotFound:
		httpStatus = http.StatusNotFound
		wsStatus = wsStatusNormalClosure

	case database.ErrGameNotFound:
		httpStatus = http.StatusNotFound
		wsStatus = wsStatusNormalClosure
	case database.ErrMaxGames:
		httpStatus = http.StatusForbidden
		wsStatus = wsStatusPolicyViolation

	default:
		httpStatus = http.StatusInternalServerError
		wsStatus = wsStatusInternalError
		res = err
	}

	api.err(writer, res, httpStatus, wsStatus)
}

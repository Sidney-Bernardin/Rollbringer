package handler

import (
	"io"
	"net/http"
	"rollbringer/pkg/domain"

	"github.com/pkg/errors"
	"golang.org/x/net/websocket"
)

const (
	wsStatusNormalClosure   = 1000
	wsStatusUnsupportedData = 1003
	wsStatusPolicyViolation = 1008
	wsStatusInternalError   = 1011
)

func (h *Handler) err(writer io.Writer, e error, httpStatus, wsStatus int) {
	switch w := writer.(type) {
	case http.ResponseWriter:

		if httpStatus >= 500 {
			h.Logger.Error().Stack().Err(e).Msg("Internal server error")
			e = errors.New("internal server error")
		}

		http.Error(w, e.Error(), httpStatus)

	case *websocket.Conn:

		if wsStatus == 1011 {
			h.Logger.Error().Stack().Err(e).Msg("Internal server error")
			e = errors.New("internal server error")
		}

		if err := w.WriteClose(wsStatus); err != nil {
			h.Logger.Error().Stack().Err(err).Msg("Cannot write close status")
		}

		w.Close()
	}
}

func (h *Handler) domainErr(writer io.Writer, err error) {

	var (
		res = errors.Cause(err)

		httpStatus int
		wsStatus   int
	)

	switch res {
	case domain.ErrUnauthorized:
		httpStatus = http.StatusUnauthorized
		wsStatus = wsStatusPolicyViolation

	case domain.ErrUserNotFound:
		httpStatus = http.StatusNotFound
		wsStatus = wsStatusNormalClosure

	case domain.ErrGameNotFound:
		httpStatus = http.StatusNotFound
		wsStatus = wsStatusNormalClosure
	case domain.ErrMaxGames:
		httpStatus = http.StatusForbidden
		wsStatus = wsStatusPolicyViolation

	case domain.ErrPlayMaterialNotFound:
		httpStatus = http.StatusNotFound
		wsStatus = wsStatusNormalClosure

	default:
		httpStatus = http.StatusInternalServerError
		wsStatus = wsStatusInternalError
		res = err
	}

	h.err(writer, res, httpStatus, wsStatus)
}

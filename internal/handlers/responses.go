package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"

	"rollbringer/internal"
)

var problemDetailStatusCodes = map[internal.PDType]int{
	internal.PDTypeServerError:  http.StatusInternalServerError,
	internal.PDTypeUnauthorized: http.StatusUnauthorized,
	internal.PDTypeInvalidID:    http.StatusBadRequest,
	internal.PDTypeInvalidView:  http.StatusBadRequest,
	internal.PDTypeInvalidEvent: http.StatusBadRequest,
	internal.PDTypeInvalidJSON:  http.StatusBadRequest,

	internal.PDTypeUserNotFound:    http.StatusNotFound,
	internal.PDTypeSessionNotFound: http.StatusNotFound,

	internal.PDTypeMaxGames:     http.StatusForbidden,
	internal.PDTypeGameNotFound: http.StatusNotFound,

	internal.PDTypePDFNotFound:          http.StatusNotFound,
	internal.PDTypeInvalidPDFName:       http.StatusBadRequest,
	internal.PDTypeInvalidPDFPageNumber: http.StatusBadRequest,
	internal.PDTypeInvalidPDFFieldName:  http.StatusBadRequest,
	internal.PDTypeNotSubscribedToPDF:   http.StatusForbidden,

	internal.PDTypeInvalidDie: http.StatusBadRequest,
}

func (h *BaseHandler) Err(w io.Writer, r *http.Request, err error) {
	pd := internal.HandleError(r.Context(), h.Logger, err)

	switch w.(type) {
	case *websocket.Conn:
		h.Render(w, r, 0, internal.EventWrapper[any]{
			Event:   internal.EventError,
			Payload: pd,
		})

	case http.ResponseWriter:
		h.Render(w, r, problemDetailStatusCodes[pd.Type], pd)
	}
}

func (h *BaseHandler) Render(w io.Writer, r *http.Request, httpStatusCode int, data any) {
	if rw, ok := w.(http.ResponseWriter); ok {
		rw.WriteHeader(httpStatusCode)
	}

	switch res := data.(type) {
	case templ.Component:
		if err := res.Render(r.Context(), w); err != nil {
			internal.HandleError(r.Context(), h.Logger, errors.Wrap(err, "cannot render component"))
			return
		}

	default:
		b, err := json.Marshal(res)
		if err != nil {
			internal.HandleError(r.Context(), h.Logger, errors.Wrap(err, "cannot encode JSON"))
			return
		}

		if _, err := w.Write(b); err != nil {
			internal.HandleError(r.Context(), h.Logger, errors.Wrap(err, "cannot write response"))
			return
		}
	}
}

package handlers

import (
	"io"
	"net/http"

	"github.com/a-h/templ"

	"rollbringer/internal"
	"rollbringer/internal/views"
)

var problemDetailStatusCodes = map[internal.PDType]int{
	internal.PDTypeServerError:  http.StatusInternalServerError,
	internal.PDTypeUnauthorized: http.StatusUnauthorized,
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
	
	var httpStatusCode int
	if _, ok := w.(http.ResponseWriter); ok {
		httpStatusCode = problemDetailStatusCodes[pd.Type]
	}

	h.Render(w, r, httpStatusCode, views.ProblemDetail(pd))
}

func (h *BaseHandler) Render(w io.Writer, r *http.Request, httpStatusCode int, data templ.Component) {
	if rw, ok := w.(http.ResponseWriter); ok {
		rw.WriteHeader(httpStatusCode)
	}

	if err := data.Render(r.Context(), w); err != nil {
		h.Logger.Error("Cannot render component", "err", err.Error())
	}
}

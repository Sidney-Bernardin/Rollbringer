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
	internal.PDTypeInvalidEvent: http.StatusUnprocessableEntity,

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

func (h *Handler) Err(w http.ResponseWriter, r *http.Request, err error) {
	pd := h.svc.HandleError(r.Context(), err)
	h.Render(w, r, problemDetailStatusCodes[pd.Type], views.ProblemDetail(pd))
}

func (h *Handler) Render(w io.Writer, r *http.Request, statusCode int, data templ.Component) {
	if rw, ok := w.(http.ResponseWriter); ok {
		rw.WriteHeader(statusCode)
	}

	if err := data.Render(r.Context(), w); err != nil {
		h.Logger.Error("Cannot render component", "err", err.Error())
	}
}

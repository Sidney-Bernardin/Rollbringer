package handler

import (
	"context"
	"io"
	"net/http"
	"rollbringer/pkg/domain"
	"rollbringer/pkg/views"

	"github.com/a-h/templ"
	"github.com/pkg/errors"
)

var normalErrorStatusCodes = map[domain.NEType]int{
	domain.NETypeServerError:           http.StatusInternalServerError,
	domain.NETypeCannotDecodeRequest:   http.StatusUnprocessableEntity,
	domain.NETypeUnauthorized:          http.StatusUnauthorized,
	domain.NETypeInvalidView:           http.StatusBadRequest,
	domain.NETypeInvalidEventOperation: http.StatusBadRequest,

	domain.NETypeUserNotFound:    http.StatusNotFound,
	domain.NETypeSessionNotFound: http.StatusNotFound,

	domain.NETypeMaxGames:     http.StatusForbidden,
	domain.NETypeGameNotFound: http.StatusNotFound,

	domain.NETypePDFNotFound:          http.StatusNotFound,
	domain.NETypeInvalidPDFName:       http.StatusBadRequest,
	domain.NETypeInvalidPDFPageNumber: http.StatusBadRequest,
	domain.NETypeInvalidPDFFieldName:  http.StatusBadRequest,
	domain.NETypeNotSubscribedToPDF:   http.StatusForbidden,
}

func (h *Handler) handleError(ctx context.Context, err error) *domain.NormalError {

	normalErr, ok := errors.Cause(err).(*domain.NormalError)
	if !ok {
		h.Logger.Error().Stack().Err(err).Msg("Server error")

		instance, _ := ctx.Value(domain.CtxKeyInstance).(string)
		normalErr = &domain.NormalError{
			Type:     domain.NETypeServerError,
			Instance: instance,
		}
	}

	return normalErr
}

func (h *Handler) err(w http.ResponseWriter, r *http.Request, err error) {
	normalErr := h.handleError(r.Context(), err)

	statusCode, ok := normalErrorStatusCodes[normalErr.Type]
	if !ok {
		h.Logger.Error().Stack().Err(err).Msg("Received problem-detail with unknown type")
	}

	h.render(w, r, statusCode, views.NormalError(normalErr))
}

func (h *Handler) render(w io.Writer, r *http.Request, statusCode int, data templ.Component) {

	if rw, ok := w.(http.ResponseWriter); ok {
		rw.WriteHeader(statusCode)
	}

	if err := data.Render(r.Context(), w); err != nil {
		h.Logger.Error().Stack().Err(err).Msg("Cannot render component")
	}
}

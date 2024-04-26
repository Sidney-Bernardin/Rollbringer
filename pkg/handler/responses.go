package handler

import (
	"io"
	"net/http"
	"rollbringer/pkg/domain"
	"rollbringer/pkg/views/components"

	"github.com/a-h/templ"
	"github.com/pkg/errors"
)

func (h *Handler) err(w http.ResponseWriter, r *http.Request, err error) {

	statusCode := http.StatusInternalServerError

	if pd, ok := errors.Cause(err).(*domain.ProblemDetail); ok {

		switch pd.Type {
		case domain.PDTypeUnauthorized:
			statusCode = http.StatusUnauthorized

		case domain.PDTypeUserNotFound:
			statusCode = http.StatusNotFound

		case domain.PDTypeGameNotFound:
			statusCode = http.StatusNotFound

		case domain.PDTypeMaxGames:
			statusCode = http.StatusForbidden

		case domain.PDTypePDFNotFound:
			statusCode = http.StatusNotFound

		case domain.PDTypeInvalidPDFName:
			statusCode = http.StatusBadRequest

		default:
			h.Logger.Error().Stack().Err(err).Msg("Received problem-detail with unknown type")
		}
	}

	h.renderErr(w, r, statusCode, err)
}

func (h *Handler) renderErr(w io.Writer, r *http.Request, statusCode int, err error) {

	pd, ok := errors.Cause(err).(*domain.ProblemDetail)
	if !ok {
		h.Logger.Error().Stack().Err(err).Msg("Server error")
		pd = &domain.ProblemDetail{
			Type: domain.PDTypeServerError,
		}
	}

	h.render(w, r, statusCode, components.Error(pd))
}

func (h *Handler) render(w io.Writer, r *http.Request, statusCode int, data templ.Component) {

	if rw, ok := w.(http.ResponseWriter); ok {
		rw.WriteHeader(statusCode)
	}

	if err := data.Render(r.Context(), w); err != nil {
		h.Logger.Error().Stack().Err(err).Msg("Cannot render component")
	}
}

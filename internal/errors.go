package internal

import (
	"context"
	"log/slog"

	"github.com/pkg/errors"
)

type PDType string

const CtxKeyInstance CtxKey = "instance"

const (
	PDTypeServerError  PDType = "server_error"
	PDTypeUnauthorized PDType = "unauthorized"

	PDTypeInvalidID    PDType = "invalid_id"
	PDTypeInvalidView  PDType = "invalid_view"
	PDTypeInvalidEvent PDType = "invalid_event"
	PDTypeInvalidJSON  PDType = "invalid_json"

	PDTypeUserNotFound    PDType = "user_not_found"
	PDTypeSessionNotFound PDType = "session_not_found"

	PDTypeMaxGames     PDType = "max_games_reached"
	PDTypeGameNotFound PDType = "game_not_found"

	PDTypePDFNotFound          PDType = "pdf_not_found"
	PDTypeInvalidPDFName       PDType = "invalid_pdf_name"
	PDTypeInvalidPDFPageNumber PDType = "invalid_pdf_page_number"
	PDTypeInvalidPDFFieldName  PDType = "invalid_pdf_field_name"
	PDTypeNotSubscribedToPDF   PDType = "not_subscribed_to_pdf"

	PDTypeInvalidDie         PDType = "invalid_die"
	PDTypeInvalidChatMessage PDType = "invalid_chat_message"
)

type PDOpts struct {
	Type   PDType
	Detail string
	Extra  map[string]any
}

func NewProblemDetail(ctx context.Context, opts PDOpts) *ProblemDetail {
	instance, _ := ctx.Value(CtxKeyInstance).(string)

	return &ProblemDetail{
		Instance: instance,
		Type:     opts.Type,
		Detail:   opts.Detail,
		Extra:    opts.Extra,
	}
}

// IsDetailed checks if the error is a ProblemDetail and has the PDType.
func IsDetailed(err error, t PDType) bool {
	pd, ok := errors.Cause(err).(*ProblemDetail)
	return ok && pd.Type == t
}

func HandleError(ctx context.Context, logger *slog.Logger, err error) *ProblemDetail {
	pd, ok := errors.Cause(err).(*ProblemDetail)
	if !ok {
		logger.Error("Server error", "err", err.Error())
		pd = NewProblemDetail(ctx, PDOpts{
			Type: PDTypeServerError,
		})
	}
	return pd
}

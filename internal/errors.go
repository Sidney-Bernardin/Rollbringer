package internal

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type PDType string

const CtxKeyInstance CtxKey = "instance"

const (
	PDTypeServerError  PDType = "server_error"
	PDTypeUnauthorized PDType = "unauthorized"
	PDTypeInvalidView  PDType = "invalid_view"
	PDTypeInvalidEvent PDType = "invalid_event"

	PDTypeUserNotFound    PDType = "user_not_found"
	PDTypeSessionNotFound PDType = "session_not_found"

	PDTypeMaxGames     PDType = "max_games_reached"
	PDTypeGameNotFound PDType = "game_not_found"

	PDTypePDFNotFound          PDType = "pdf_not_found"
	PDTypeInvalidPDFName       PDType = "invalid_pdf_name"
	PDTypeInvalidPDFPageNumber PDType = "invalid_pdf_page_number"
	PDTypeInvalidPDFFieldName  PDType = "invalid_pdf_field_name"
	PDTypeNotSubscribedToPDF   PDType = "not_subscribed_to_pdf"

	PDTypeInvalidDie PDType = "invalid_die"
)

type ProblemDetail struct {
	Instance string `json:"instance,omitempty"`
	Type     PDType `json:"type"`
	Detail   string `json:"detail,omitempty"`
	Extra    map[string]any
}

type PDOpts struct {
	Type   PDType
	Detail string
	Extra  map[string]any
}

func NewProblemDetail(ctx context.Context, opts PDOpts) *ProblemDetail {
	return &ProblemDetail{
		Instance: ctx.Value(CtxKeyInstance).(string),
		Type:     opts.Type,
		Detail:   opts.Detail,
		Extra:    opts.Extra,
	}
}

func (pd *ProblemDetail) Error() string {
	return fmt.Sprintf("%s: %s", pd.Type, pd.Detail)
}

// IsDetailed checks if the error is a ProblemDetail and has the PDType.
func IsDetailed(err error, t PDType) bool {
	pd, ok := errors.Cause(err).(*ProblemDetail)
	return ok && pd.Type == t
}

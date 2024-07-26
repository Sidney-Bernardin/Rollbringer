package domain

import (
	"github.com/pkg/errors"
)

type NEType string

const (
	NETypeServerError           NEType = "server_error"
	NETypeCannotDecodeRequest   NEType = "cannot_decode_event"
	NETypeUnauthorized          NEType = "unauthorized"
	NETypeInvalidView           NEType = "invalid_view"
	NETypeInvalidEventOperation NEType = "invalid_event_type"

	NETypeUserNotFound    NEType = "user_not_found"
	NETypeSessionNotFound NEType = "session_not_found"

	NETypeMaxGames     NEType = "max_games_reached"
	NETypeGameNotFound NEType = "game_not_found"

	NETypePDFNotFound          NEType = "pdf_not_found"
	NETypeInvalidPDFName       NEType = "invalid_pdf_name"
	NETypeInvalidPDFPageNumber NEType = "invalid_pdf_page_number"
	NETypeInvalidPDFFieldName  NEType = "invalid_pdf_field_name"
	NETypeNotSubscribedToPDF   NEType = "not_subscribed_to_pdf"

	NETypeInvalidDie NEType = "invalid_die"
)

type NormalError struct {
	Instance string `json:"instance,omitempty"`
	Type     NEType `json:"type"`
	Detail   string `json:"detail,omitempty"`
}

func (pd *NormalError) Error() string {
	return string(pd.Type)
}

// IsNormal checks if the error is a NormalError and has the NEType.
func IsNormal(err error, t NEType) bool {
	pd, ok := errors.Cause(err).(*NormalError)
	return ok && pd.Type == t
}

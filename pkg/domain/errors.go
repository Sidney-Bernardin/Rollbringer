package domain

import "github.com/pkg/errors"

type pdType string

const (
	PDTypeServerError           pdType = "server error"
	PDTypeCannotDecodeEvent     pdType = "cannot decode event"
	PDTypeUnauthorized          pdType = "unauthorized"
	PDTypeUnknownView           pdType = "unknown view"
	PDTypeUnknownEventOperation pdType = "unknown event operation"

	PDTypeUserNotFound    pdType = "user not found"
	PDTypeSessionNotFound pdType = "session not found"

	PDTypeMaxGames     pdType = "max games reached"
	PDTypeGameNotFound pdType = "game not found"

	PDTypePDFNotFound          pdType = "pdf not found"
	PDTypePDFPageNotFound      pdType = "pdf page not found"
	PDTypeInvalidPDFName       pdType = "invalid pdf name"
	PDTypeInvalidPDFFieldName  pdType = "invalid pdf field name"
	PDTypeInvalidPDFPageNumber pdType = "invalid pdf page number"
	PDTypeNotSubscribedToPDF   pdType = "invalid pdf field value"
)

type ProblemDetail struct {
	Type   pdType `json:"type"`
	Detail string `json:"detail,omitempty"`
}

func (pd *ProblemDetail) Error() string {
	return string(pd.Type)
}

// IsProblemDetail checks if the error is a problem-detail with the pdType.
func IsProblemDetail(err error, t pdType) bool {
	pd, ok := errors.Cause(err).(*ProblemDetail)
	return ok && pd.Type == t
}

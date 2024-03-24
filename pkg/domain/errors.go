package domain

import "github.com/pkg/errors"

type pdType string

const (
	PDTypeServerError           pdType = "server error"
	PDTypeCannotDecodeEvent     pdType = "cannot decode event"
	PDTypeUnauthorized          pdType = "unauthorized"
	PDTypeInvalidEventOperation pdType = "invalid event operation"

	PDTypeUserNotFound pdType = "user not found"

	PDTypeMaxGames     pdType = "max games reached"
	PDTypeGameNotFound pdType = "game not found"

	PDTypePDFNotFound          pdType = "pdf not found"
	PDTypeInvalidPDFName       pdType = "invalid pdf name"
	PDTypeInvalidPDFFieldValue pdType = "invalid pdf field value"
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

package server

import (
	"fmt"

	googleUUID "github.com/google/uuid"
)

type UserError struct {
	Type    UserErrorType  `json:"type"`
	Message string         `json:"message,omitempty"`
	Details map[string]any `json:"details,omitempty"`
}

func (err *UserError) Error() string {
	return fmt.Sprintf("%s: %s %v", err.Type, err.Message, err.Details)
}

type UserErrorType string

const (
	UserErrorTypeInternalServerError     UserErrorType = "internal-server-error"
	UserErrorTypeUnauthorized            UserErrorType = "unauthorized"
	UserErrorTypeUUIDInvalid             UserErrorType = "uuid-invalid"
	UserErrorTypeGoogleUserAlreadyExists UserErrorType = "google-user-already-exists"
)

/////

type UUID struct {
	googleUUID.UUID
}

func NewRandomUUID() UUID {
	return UUID{googleUUID.Must(googleUUID.NewRandom())}
}

func ParseUUID(uuid string) (ret UUID, err error) {
	gUUID, err := googleUUID.Parse(uuid)
	if err != nil {
		return ret, &UserError{
			Type:    UserErrorTypeUUIDInvalid,
			Message: err.Error(),
		}
	}
	return UUID{gUUID}, nil
}

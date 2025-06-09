package server

import (
	"fmt"
	"strings"

	googleUUID "github.com/google/uuid"
)

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

func (uuid UUID) Escape() string {
	return strings.Replace(uuid.String(), `-`, `\-`, -1)
}

/////

type UserError struct {
	Type    UserErrorType  `json:"type"`
	Message string         `json:"message,omitempty"`
	Details map[string]any `json:"details,omitempty"`
}

func NewUserError(typ UserErrorType, message string, details map[string]any) *UserError {
	return &UserError{typ, message, details}
}

func (err *UserError) Error() string {
	return fmt.Sprintf("%s: %s %v", err.Type, err.Message, err.Details)
}

type UserErrorType string

const (
	UserErrorTypeInternalServerError UserErrorType = "internal-server-error"
	UserErrorTypeUnauthorized        UserErrorType = "unauthorized"
	UserErrorTypeUUIDInvalid         UserErrorType = "uuid-invalid"

	UserErrorTypeUsernameInvalid         UserErrorType = "username-invalid"
	UserErrorTypeUsernameTaken           UserErrorType = "username-taken"
	UserErrorTypePasswordInvalid         UserErrorType = "password-invalid"
	UserErrorTypeGoogleUserAlreadyExists UserErrorType = "google-user-already-exists"
	UserErrorTypeGoogleUserNotExists     UserErrorType = "google-user-not-exists"
	UserErrorTypeUserNotFound            UserErrorType = "user-not-found"

	UserErrorTypeRoomNameInvalid UserErrorType = "room-name-invalid"
	UserErrorTypeRoomNotFound    UserErrorType = "room-not-found"
)

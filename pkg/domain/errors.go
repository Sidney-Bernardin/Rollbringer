package domain

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

/////

type detailedError struct {
	msg   string
	attrs map[string]any
	child error
}

func Wrap(err error, msg string, attrs map[string]any) error {
	if err == nil {
		return nil
	}

	return &detailedError{
		msg:   msg,
		attrs: attrs,
		child: err,
	}
}

func (err *detailedError) Error() string {
	return fmt.Sprintf("%s: %v", err.msg, err.child)
}

func (err *detailedError) Unwrap() error {
	return err.child
}

func (err *detailedError) Cause() error {
	return err.child
}

/////

type UserError struct {
	Instance    string         `json:"instance,omitempty"`
	Type        UserErrorType  `json:"type"`
	Description string         `json:"description,omitempty"`
	Attrs       map[string]any `json:"attrs,omitempty"`
}

type UserErrorType string

const (
	UsrErrTypeServerError          UserErrorType = "server_error"
	UsrErrTypeCannotProcessRequest UserErrorType = "cannot_process_request"
	UsrErrTypeUnauthorized         UserErrorType = "unauthorized"
	UsrErrTypeRecordNotFound       UserErrorType = "record_not_found"

	UsrErrTypeGoogleUserDoesNotExists UserErrorType = "google_user_does_not_exist"
	UsrErrTypeGoogleUserAlreadyExists UserErrorType = "google_user_already_exists"

	UsrErrTypeSpotifyUserDoesNotExists UserErrorType = "spotify_user_does_not_exist"
	UsrErrTypeSpotifyUserAlreadyExists UserErrorType = "spotify_user_already_exists"
)

func UserErr(ctx context.Context, t UserErrorType, desc string, attrs map[string]any) *UserError {
	instance, _ := ctx.Value("instance").(string)
	return &UserError{
		Instance:    instance,
		Type:        t,
		Description: desc,
		Attrs:       attrs,
	}
}

func (err *UserError) Error() string {
	return fmt.Sprintf("%s: %s: %s", err.Instance, err.Type, err.Description)
}

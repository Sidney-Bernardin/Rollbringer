package domain

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrNotFound = errors.New("not found")
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
	return fmt.Sprintf("%s: %s", err.msg, err.child.Error())
}

func (err *detailedError) Unwrap() error {
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

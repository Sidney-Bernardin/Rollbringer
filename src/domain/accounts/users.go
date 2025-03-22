package accounts

import (
	"context"

	"github.com/pkg/errors"

	"rollbringer/src/domain"
)

type Username string

func ParseUsername(str string) (Username, error) {
	if len(str) == 0 || 25 < len(str) {
		return "", &domain.DomainError{
			Type:        DomainErrorTypeUsernameInvalid,
			Description: "Must be between 1 and 25 characters",
			Details:     map[string]any{"username": str},
		}
	}

	return Username(str), nil
}

/////

type CmdUserCreate struct {
	Username Username
}

type ArgsUserCreate struct {
	Username string
}

func (svc *service) UserCreate(ctx context.Context, view any, args *ArgsUserCreate) (err error) {
	var cmd CmdUserCreate

	cmd.Username, err = ParseUsername(args.Username)
	if err != nil {
		return errors.Wrap(err, "cannot parse username")
	}

	if err := svc.db.UserCreate(ctx, view, &cmd); err != nil {
		if errors.Is(err, domain.ErrEntityConflict) {
			return &domain.DomainError{
				Type:        DomainErrorTypeUsernameTaken,
				Description: "A user with the given username already exists.",
				Details:     map[string]any{"username": cmd.Username},
			}
		}

		return errors.Wrap(err, "cannot insert room")
	}

	return nil
}

/////

func (svc *service) UserGetByUsername(ctx context.Context, view any, usernameStr string) error {

	username, err := ParseUsername(usernameStr)
	if err != nil {
		return errors.Wrap(err, "cannot parse username")
	}

	if err := svc.db.UserGetByUsername(ctx, view, username); err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &domain.DomainError{
				Type:        domain.DomainErrorTypeEntityNotFound,
				Description: "Cannot find a user with the given username",
				Details:     map[string]any{"username": username},
			}
		}

		return errors.Wrap(err, "cannot get user by username")
	}

	return nil
}

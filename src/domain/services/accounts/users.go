package accounts

import (
	"context"
	"rollbringer/src/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type User struct {
	ID uuid.UUID

	GoogleID  *string
	SpotifyID *string

	Username       Username
	ProfilePicture string
}

type Username string

const ExternalErrorTypeInvalidUsername domain.ExternalErrorType = "invalid-username"

func ParseUsername(str string) (Username, error) {
	if len(str) == 0 || 25 < len(str) {
		return "", &domain.ExternalError{
			Type:    ExternalErrorTypeInvalidUsername,
			Msg:     "Must be between 1 and 25 characters",
			Details: map[string]any{"username": str},
		}
	}

	return Username(str), nil
}

func (svc *service) GetUserByUserID(ctx context.Context, userID uuid.UUID) (*User, error) {
	user, err := svc.db.GetUserByUserID(ctx, userID)
	return user, errors.Wrap(err, "cannot get user by user-ID")
}

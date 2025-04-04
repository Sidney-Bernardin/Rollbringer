package accounts

import (
	"context"
	"rollbringer/src"
	"rollbringer/src/domain"

	"github.com/pkg/errors"
)

type CSRFToken string

func NewCSRFToken() CSRFToken {
	return CSRFToken(src.CreateRandomString())
}

/////

func (svc *service) Auth(ctx context.Context, sessionIDStr string, csrfToken *string) (*ViewSessionInfo, error) {

	sessionID, err := domain.ParseUUID(sessionIDStr)
	if err != nil {
		return nil, nil
	}

	var sessionInfo ViewSessionInfo
	if csrfToken == nil {
		err = svc.db.SessionGetByID(ctx, &sessionInfo, sessionID)
		err = errors.Wrap(err, "cannot get session by ID")
	} else {
		err = svc.db.SessionGetByIDAndCSRFToken(ctx, &sessionInfo, sessionID, CSRFToken(*csrfToken))
		err = errors.Wrap(err, "cannot get session by ID and CSRF-token")
	}

	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "cannot get session")
	}

	return &sessionInfo, nil
}

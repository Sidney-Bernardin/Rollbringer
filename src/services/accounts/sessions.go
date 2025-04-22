package accounts

import (
	"context"

	"github.com/pkg/errors"

	"rollbringer/src"
	"rollbringer/src/services/accounts/models"
)

func (svc *service) Auth(ctx context.Context, sessionIDStr string, csrfToken *string) (*models.Session, error) {

	// Parse the sessionID.
	sessionID, err := src.ParseUUID(sessionIDStr)
	if err != nil {
		return nil, nil
	}

	// Get the session.
	var session *models.Session
	if csrfToken == nil {
		session, err = svc.database.GetSessionByID(ctx, sessionID)
	} else {
		session, err = svc.database.GetSessionByIDAndCSRFToken(ctx, sessionID, models.CSRFToken(*csrfToken))
	}

	if err != nil {
		if errors.Is(err, src.ErrEntityNotFound) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "database cannot get session")
	}

	return session, nil
}

package commands

import (
	"rollbringer/server/domain"
	"rollbringer/server/domain/play/results"

	"github.com/google/uuid"
)

type UUID uuid.UUID

func NewUUID(uuidStr string) (UUID, error) {
	uuidParsed, err := uuid.Parse(uuidStr)
	if err != nil {
		return UUID(uuid.Nil), &domain.DomainError{
			Type:        results.DomainErrorTypeUUIDInvalid,
			Description: err.Error(),
			Details:     map[string]any{"uuid": uuidStr},
		}
	}

	return UUID(uuidParsed), nil
}

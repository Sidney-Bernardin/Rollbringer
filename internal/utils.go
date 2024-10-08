package internal

import (
	"context"

	"github.com/google/uuid"
)

type CtxKey string

var PDFSchemaPageNames = map[string][]string{
	"DND_CHARACTER_SHEET": {"main", "info", "spells"},
}

func OptionalID(ctx context.Context, s string) (*uuid.UUID, error) {
	if s == "" {
		return nil, nil
	}

	id, err := uuid.Parse(s)
	if err != nil {
		return nil, NewProblemDetail(ctx, PDOpts{
			Type:   PDTypeInvalidID,
			Detail: err.Error(),
		})
	}

	return &id, nil
}

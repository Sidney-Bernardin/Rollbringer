package internal

import "github.com/google/uuid"

type CtxKey string

var PDFSchemaPageNames = map[string][]string{
	"DND_CHARACTER_SHEET": {"main", "info", "spells"},
}

func OptionalUUID(s string) *uuid.UUID {
	if s == "" {
		return nil
	}

	parsed, _ := uuid.Parse(s)
	return &parsed
}

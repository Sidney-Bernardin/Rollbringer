package domain

import "github.com/google/uuid"

func ParseUUIDs(ids ...*string) {
	for _, id := range ids {
		parsed, _ := uuid.Parse(*id)
		*id = parsed.String()
	}
}

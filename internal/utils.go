package internal

import (
	"context"
	"strings"
)

type CtxKey string

var PDFSchemaPageNames = map[string][]string{
	"DND_CHARACTER_SHEET": {"main", "info", "spells"},
}

func ParseViewQuery[T UserView | SessionView | GameView | PDFView](ctx context.Context, viewQuery string) (map[string]T, error) {
	ret := map[string]T{}

	for _, view := range strings.Split(viewQuery, ",") {
		if field := strings.Split(view, "-"); len(field) == 2 {
			ret[field[0]] = T(field[1])
			continue
		}

		return nil, NewProblemDetail(ctx, PDOpts{
			Type: PDTypeInvalidView,
		})
	}

	return ret, nil
}

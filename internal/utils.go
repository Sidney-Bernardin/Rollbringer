package internal

import "strings"

type CtxKey string

var PDFSchemaPageNames = map[string][]string{
	"DND_CHARACTER_SHEET": {"main", "info", "spells"},
}

func ParseViews(views string) {
	// pdf-all,owner-name,game-all
	for _, view := range strings.Split(views, ",") {
		viewParts := strings.Split(view, "-")

		switch viewParts[0] {
		case GameViews[0]:
		}
	}
}

package internal

type CtxKey string

var PDFSchemaPageNames = map[string][]string{
	"DND_CHARACTER_SHEET": {"main", "info", "spells"},
}

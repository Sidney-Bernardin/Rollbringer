package domain

type ctxKey string

const (
	CtxKeyInstance ctxKey = "instance"
)

var PDFSchemaPageNames = map[string][]string{
	"DND_CHARACTER_SHEET": {"main", "info", "spells"},
}

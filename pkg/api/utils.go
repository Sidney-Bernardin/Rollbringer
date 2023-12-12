package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"html/template"
)

var templateUtils = template.FuncMap{
	"mustJSONIFY": func(v any) string {
		b, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		return string(b)
	},
}

func mustGetRandHexStr() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

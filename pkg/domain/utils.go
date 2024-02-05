package domain

import (
	"encoding/json"
	"unicode"
)

var PDFSchemaPageCount = map[string]int{
	"DND_CHARACTER_SHEET": 3,
}

func DecodeGameEvent(e []byte) *GameEvent {

	var eventModel GameEvent
	if err := json.Unmarshal(e, &eventModel); err != nil {
		return nil
	}

	var eventMap map[string]any
	if err := json.Unmarshal(e, &eventMap); err != nil {
		return nil
	}

	delete(eventMap, "HEADERS")
	delete(eventMap, "TYPE")

	eventModel.PDFFields = make(map[string]string)
	for k, v := range eventMap {
		if unicode.IsUpper(rune(k[0])) {
			if valStr, ok := v.(string); ok {
				eventModel.PDFFields[k] = valStr
			}
		}
	}

	return &eventModel
}

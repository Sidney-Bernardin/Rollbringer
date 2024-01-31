package domain

import (
	"encoding/json"
)

func DecodeGameEvent(b []byte) *GameEvent {

	var event map[string]any
	if err := json.Unmarshal(b, &event); err != nil {
		return nil
	}

	var ret GameEvent
	if err := json.Unmarshal(b, &ret); err != nil {
		return nil
	}

	for _, v := range [2]string{"HEADERS", "TYPE"} {
		delete(event, v)
	}

	ret.Body = make(map[string]string)
	for k, v := range event {
		ret.Body[k] = v.(string)
	}

	return &ret
}

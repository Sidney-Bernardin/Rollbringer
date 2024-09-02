package internal

import (
	"context"
	"encoding/json"
	"fmt"
)

type CtxKey string

var PDFSchemaPageNames = map[string][]string{
	"DND_CHARACTER_SHEET": {"main", "info", "spells"},
}

func JSONEncodeEvent(ctx context.Context, event Event) ([]byte, error) {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return nil, NewProblemDetail(ctx, PDOpts{
			Type:   PDTypeInvalidEvent,
			Detail: err.Error(),
		})
	}
	return eventBytes, nil
}

func JSONDecodeEvent(ctx context.Context, eventBytes []byte) (Event, error) {

	var baseEvent BaseEvent
	if err := json.Unmarshal(eventBytes, &baseEvent); err != nil {
		return nil, NewProblemDetail(ctx, PDOpts{
			Type:   PDTypeInvalidEvent,
			Detail: err.Error(),
		})
	}

	event, ok := eventTypes[baseEvent.Type]
	if !ok {
		return nil, NewProblemDetail(ctx, PDOpts{
			Type:   PDTypeInvalidEvent,
			Detail: fmt.Sprintf("'%s' is an invalid event-type.", baseEvent.Type),
			Extra: map[string]any{
				"event_type": baseEvent.Type,
			},
		})
	}

	if err := json.Unmarshal(eventBytes, &event); err != nil {
		return nil, NewProblemDetail(ctx, PDOpts{
			Type:   PDTypeInvalidEvent,
			Detail: err.Error(),
		})
	}

	return event, nil
}

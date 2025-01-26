package domain

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/pkg/errors"
)

type DefinePayloadFunc func(Operation) (payload any)

func DecodeEvent(bytes []byte, dpFunc DefinePayloadFunc) (*Event, error) {
	op := struct {
		Operation Operation `json:"operation"`
	}{}

	// Decode the event's operation.
	if err := json.Unmarshal(bytes, &op); err != nil {
		return nil, Wrap(err, "cannot JSON decode event operation", nil)
	}

	event := &Event{
		Operation: op.Operation,
		Payload:   dpFunc(op.Operation),
	}

	// Decode the event.
	if err := json.Unmarshal(bytes, event); err != nil {
		return event, Wrap(err, "cannot JSON decode event", nil)
	}

	return event, nil
}

// HandleError checks if the error's cause is a *UserError. If so, the cause is
// returned. If not, it's treated as a server error. This means logging the error
// with any of it's wrapped attributes and returning a new *UserError with UsrErrTypeServerError.
func HandleError(ctx context.Context, logger *slog.Logger, level slog.Level, err error) *UserError {
	if err == nil {
		return nil
	}

	logAttrs := []slog.Attr{slog.String("err", err.Error())}

	// Loop over the error's unwrap chain.
	for e := err; e != nil; e = errors.Unwrap(e) {
		switch e := e.(type) {

		// Add the attributes to logAttrs.
		case *detailedError:
			for k, v := range e.attrs {
				logAttrs = append(logAttrs, slog.Any(k, v))
			}

		case *UserError:
			return e
		}
	}

	logger.LogAttrs(ctx, level, "Server error", logAttrs...)
	return UserErr(ctx, UsrErrTypeServerError, "", nil)
}

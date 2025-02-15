package domain

import (
	"context"
)

type PubSubRepository interface {
	Request(ctx context.Context, subject string, resPayload any, req *Event) (Operation, error)
	Publish(ctx context.Context, subject string, req *Event) error
	Subscribe(ctx context.Context, subject string, cb PubSubCallback, expectedEvents map[Operation]any) error
	Close()
}

type PubSubCallback func(context.Context, *Event) *Event

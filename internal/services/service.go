package services

import (
	"context"
	"log/slog"

	"github.com/pkg/errors"

	"rollbringer/internal"
	"rollbringer/internal/config"
)

type Servicer interface {
	HandleError(context.Context, error) *internal.ProblemDetail
}

type Service struct {
	Config *config.Config
	Logger *slog.Logger
}

func (svc *Service) HandleError(ctx context.Context, err error) *internal.ProblemDetail {
	pd, ok := errors.Cause(err).(*internal.ProblemDetail)
	if !ok {
		svc.Logger.Error("Server error", "err", err.Error())

		instance, _ := ctx.Value(internal.CtxKeyInstance).(string)
		pd = &internal.ProblemDetail{
			Type:     internal.PDTypeServerError,
			Instance: instance,
		}
	}

	return pd
}

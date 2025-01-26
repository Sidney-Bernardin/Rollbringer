package nats

import (
	"context"
	"fmt"
	"log/slog"
	"rollbringer/pkg/domain"
)

type natsLogger struct {
	logger *slog.Logger
}

func (l *natsLogger) Tracef(format string, a ...any) {
	l.logger.Log(context.Background(), domain.SlogLevelTrace, fmt.Sprintf(format, a...))
}

func (l *natsLogger) Debugf(format string, a ...any) {
	l.logger.Log(context.Background(), slog.LevelDebug, fmt.Sprintf(format, a...))
}

func (l *natsLogger) Noticef(format string, a ...any) {
	l.logger.Log(context.Background(), slog.LevelInfo, fmt.Sprintf(format, a...))
}

func (l *natsLogger) Warnf(format string, a ...any) {
	l.logger.Log(context.Background(), slog.LevelWarn, fmt.Sprintf(format, a...))
}

func (l *natsLogger) Errorf(format string, a ...any) {
	l.logger.Log(context.Background(), slog.LevelError, fmt.Sprintf(format, a...))
}

func (l *natsLogger) Fatalf(format string, a ...any) {
	l.logger.Log(context.Background(), domain.SlogLevelFatal, fmt.Sprintf(format, a...))
}

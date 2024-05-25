package mlog

import (
	"context"
	"log/slog"
)

func L(ctx context.Context) *slog.Logger {
	logger := ctx.Value(loggerKey)
	switch logger := logger.(type) {
	case *slog.Logger:
		return logger
	default:
		return slog.Default()
	}
}

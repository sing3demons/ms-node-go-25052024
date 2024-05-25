package mlog

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type contextKey string

const (
	sessionKey contextKey = "session"
	loggerKey  contextKey = "logger"
)

func logMiddleware(ctx context.Context, logger *slog.Logger) *slog.Logger {
	session, exits := ctx.Value(loggerKey).(string)
	if !exits {
		session = uuid.New().String()
	}
	logger = logger.With("session", session)
	return logger
}

func Middleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := c.GetHeader("x-request-id")
		if session == "" {
			xSession := c.GetHeader("x-session-id")
			if xSession != "" {
				session = xSession
			} else {
				session = uuid.New().String()
			}
		}
		ctx := context.WithValue(c.Request.Context(), sessionKey, session)
		l := logMiddleware(ctx, logger)
		ctx = context.WithValue(ctx, loggerKey, l)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

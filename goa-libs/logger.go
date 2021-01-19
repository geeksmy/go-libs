package goalibs

import (
	"context"

	"go.uber.org/zap"
	"goa.design/goa/v3/middleware"
)

// 日志绑定上下文
// 目前仅添加 requestID
func LoggerWithContext(ctx context.Context, logger *zap.Logger) *zap.Logger {
	requestID, ok := ctx.Value(middleware.RequestIDKey).(string)
	if ok {
		logger = logger.With(zap.String("requestID", requestID))
	}
	return logger
}

func Logger(ctx context.Context) *zap.Logger {
	return LoggerWithContext(ctx, zap.L())
}

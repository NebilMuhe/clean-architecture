package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type Logger interface {
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Fatal(ctx context.Context, msg string, fields ...zap.Field)
	With(fields ...zap.Field) *zap.Logger
	extract(ctx context.Context) []zap.Field
}

type logger struct {
	logger *zap.Logger
}

func NewLogger(log *zap.Logger) Logger {
	return &logger{
		logger: log,
	}
}

func (l *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.With(l.extract(ctx)...).Error(msg, fields...)
	// l.logger.With(l.extract(ctx)...).Error(msg,fields...)
}

func (l *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.With(l.extract(ctx)...).Info(msg, fields...)
}

func (l *logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	l.With(l.extract(ctx)...).Fatal(msg, fields...)
}

func (l *logger) With(fields ...zap.Field) *zap.Logger {
	l2 := l.logger.With(fields...)
	return l2
}

func (l *logger) extract(ctx context.Context) []zap.Field {
	var fields []zap.Field

	fields = append(fields, zap.String("time", time.Now().Format(time.RFC3339)))

	if reqID, ok := ctx.Value("x-request-id").(string); ok {
		fields = append(fields, zap.String("x-request-id", reqID))
	}

	if hitTime, ok := ctx.Value("time-since-start").(time.Time); ok {
		fields = append(fields, zap.Float64("time-since-start", float64(time.Since(hitTime).Milliseconds())))
	}

	return fields
}

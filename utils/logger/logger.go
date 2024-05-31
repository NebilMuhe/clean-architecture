package logger

import (
	"context"

	"go.uber.org/zap"
)

type Logger interface {
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Fatal(ctx context.Context, msg string, fields ...zap.Field)
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

}

func (l *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {

}

func (l *logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {

}
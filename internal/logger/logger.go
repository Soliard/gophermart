package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
)

const ctxKeyLogger ctxKey = "logger"

var Log Logger
var F = fieldBuilder{}

type ctxKey string

type Field = zap.Field

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	With(fields ...Field) Logger
}

type fieldBuilder struct{}

type zapWrapper struct {
	*zap.Logger
}

func (z zapWrapper) With(fields ...Field) Logger {
	return zapWrapper{z.Logger.With(fields...)}
}

func WithContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, ctxKeyLogger, logger)
}

func FromContext(ctx context.Context) Logger {
	if logger, ok := ctx.Value(ctxKeyLogger).(Logger); ok {
		return logger
	}
	return Log
}

func Init() error {
	z, err := zap.NewProduction()
	if err != nil {
		return err
	}
	Log = zapWrapper{z}
	return nil
}

func (fieldBuilder) String(key, value string) Field {
	return zap.String(key, value)
}

func (fieldBuilder) Int(key string, value int) Field {
	return zap.Int(key, value)
}

func (fieldBuilder) Any(key string, value interface{}) Field {
	return zap.Any(key, value)
}

func (fieldBuilder) Error(err error) Field {
	return zap.Error(err)
}

func (fieldBuilder) Duration(key string, value time.Duration) Field {
	return zap.Duration(key, value)
}

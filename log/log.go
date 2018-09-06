// Package log is a wrapper around the zap logging library with convenience methods and constants
// predefined.
package log

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap/zapcore"

	zap "go.uber.org/zap"
)

// ContextKey is a type for storing logging information inside of a context
type ContextKey string

const (
	// ContextKeyRequestHost represents the requesting host (hostname / IP)
	// Filled for every HTTP request.
	ContextKeyRequestHost = ContextKey("request_host")
	// ContextKeyEndpoint represents the requested endpoint (/v1/endpoint)
	// Filled for every HTTP request.
	ContextKeyEndpoint = ContextKey("endpoint")
	// ContextKeyRequestMethod represents the request method (GET / POST / ...)
	// Filled for every HTTP request.
	ContextKeyRequestMethod = ContextKey("request_method")
)

// Field is an alias for a zapcore.Field
type Field = zapcore.Field

var logger *zap.Logger

func init() {
	encoderConfig := zapcore.EncoderConfig{

		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "ts",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	logConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: true,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, _ = logConfig.Build(zap.AddCallerSkip(1))
}

// Printf is a wrapper around logger.Info which emulates a simple Printf command.
// This is used in places where a log.Printf-like interface is expected (swagger server).
func Printf(template string, fields ...interface{}) {
	logger.Info(fmt.Sprintf(template, fields...))
}

func Debug(ctx context.Context, msg string, fields ...Field) {
	logger.Debug(msg, append(fields, getContextFields(ctx)...)...)
}

func DebugCat(ctx context.Context, msg string, cat string, fields ...Field) {
	fields = append(fields, String("cat", cat))
	logger.Debug(msg, append(fields, getContextFields(ctx)...)...)
}

func Info(ctx context.Context, msg string, fields ...Field) {
	logger.Info(msg, append(fields, getContextFields(ctx)...)...)
}

func InfoCat(ctx context.Context, msg string, cat string, fields ...Field) {
	fields = append(fields, String("cat", cat))
	logger.Info(msg, append(fields, getContextFields(ctx)...)...)
}

func Warn(ctx context.Context, msg string, fields ...Field) {
	logger.Warn(msg, append(fields, getContextFields(ctx)...)...)
}

func WarnCat(ctx context.Context, msg string, cat string, fields ...Field) {
	fields = append(fields, String("cat", cat))
	logger.Warn(msg, append(fields, getContextFields(ctx)...)...)
}

func Error(ctx context.Context, msg string, fields ...Field) {
	logger.Error(msg, append(fields, getContextFields(ctx)...)...)
}

func ErrorCat(ctx context.Context, msg string, cat string, fields ...Field) {
	fields = append(fields, String("cat", cat))
	logger.Error(msg, append(fields, getContextFields(ctx)...)...)
}

// Field wrappers

// String is a wrapper around zap.String
func String(key string, value string) Field {
	return zap.String(key, value)
}

// Strings is a wrapper around zap.Strings
func Strings(key string, value []string) Field {
	return zap.Strings(key, value)
}

// Int is a wrapper around zap.Int
func Int(key string, value int) Field {
	return zap.Int(key, value)
}

// Bool is a wrapper around zap.Bool
func Bool(key string, value bool) Field {
	return zap.Bool(key, value)
}

// Bytes is a wrapper around zap.Bytes
func Bytes(key string, value []byte) Field {
	return zap.ByteString(key, value)
}

// Duration is a wrapper around zap.Duration
func Duration(key string, value time.Duration) Field {
	return zap.Duration(key, value)
}

func getContextFields(ctx context.Context) []zapcore.Field {
	fields := []zapcore.Field{}
	if host, ok := ctx.Value(ContextKeyRequestHost).(string); ok && host != "" {
		fields = append(fields, zap.String(string(ContextKeyRequestHost), host))
	}
	if endpoint, ok := ctx.Value(ContextKeyEndpoint).(string); ok && endpoint != "" {
		fields = append(fields, zap.String(string(ContextKeyEndpoint), endpoint))
	}
	if method, ok := ctx.Value(ContextKeyRequestMethod).(string); ok && method != "" {
		fields = append(fields, zap.String(string(ContextKeyRequestMethod), method))
	}

	return fields
}

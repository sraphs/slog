// Package log provides a global logger for log.
package slog

import (
	"context"
	"io"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

//go:generate protoc --go_out=paths=source_relative:. log.proto

// LevelLogger is a logger interface that provides logging function with levels.
type LevelLogger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Print(v ...interface{}) // Print is an alias of Info().
	Warn(v ...interface{})
	Error(v ...interface{})
	Fatal(v ...interface{})
}

// FormatLogger is a logger interface that output logs with a format.
type FormatLogger interface {
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Printf(format string, v ...interface{}) // Printf is an alias for Infof.
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}

// Control provides methods to config a logger.
type Control interface {
	SetLevel(Level) Control
	SetOutput(io.Writer) Control
}

// Kratos logger interface.
type KLogger interface {
	log.Logger
}

// FullLogger is the combination of Logger, FormatLogger, CtxLogger and Control.
type FullLogger interface {
	KLogger
	LevelLogger
	FormatLogger
	Control
	Clone() FullLogger
	WithTimestamp() FullLogger
	WithCaller() FullLogger
	WithCallerWithSkipFrameCount(skipFrameCount int) FullLogger
	WithStack() FullLogger
	WithFields(fields ...interface{}) FullLogger
}

type Level = log.Level

const (
	// LevelDebug is logger debug level.
	LevelDebug = log.LevelDebug
	// LevelInfo is logger info level.
	LevelInfo = log.LevelInfo
	// LevelWarn is logger warn level.
	LevelWarn = log.LevelWarn
	// LevelError is logger error level.
	LevelError = log.LevelError
	// LevelFatal is logger fatal level
	LevelFatal = log.LevelFatal
)

// ParseLevel takes a string level and returns the logger log level constant.
func ParseLevel(lv string) Level {
	switch strings.ToUpper(lv) {
	case "DEBUG":
		return LevelDebug
	case "INFO":
		return LevelInfo
	case "WARN":
		return LevelWarn
	case "ERROR":
		return LevelError
	case "FATAL":
		return LevelFatal
	}

	return LevelInfo
}

// loggerKey points to the value in the context where the logger is stored.
type loggerKey struct{}

// WithLogger creates a new context with the provided logger attached.
func WithLogger(ctx context.Context, logger FullLogger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// FromContext returns the logger stored in the context. If no such logger
// exists, a default logger is returned.
func FromContext(ctx context.Context) FullLogger {
	if logger, ok := ctx.Value(loggerKey{}).(FullLogger); ok {
		return logger
	}
	return DefaultLogger()
}

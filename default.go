package slog

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

var logger FullLogger

func init() {
	Init(nil)
}

func Init(c *Config) FullLogger {
	logger = New(c)
	return logger
}

func New(c *Config) FullLogger {
	var w io.Writer = os.Stdout
	if c == nil {
		c = &Config{
			Level: "info",
		}
	}

	if c.Path != "" {
		logrotate := &lumberjack.Logger{
			Filename: c.Path,
			MaxSize:  int(c.MaxSize),
			MaxAge:   int(c.MaxAge),
		}
		w = MultiLevelWriter(logrotate, os.Stdout)
	}

	l := newZerolog(w)
	lv := ParseLevel(c.Level)
	l.SetLevel(lv)
	return &defaultLogger{log: l}
}

func Clone() FullLogger {
	return logger.Clone()
}

// SetLevel sets the current global log level.
func SetLevel(lv Level) Control {
	return logger.SetLevel(lv)
}

// SetOutput sets the global logger output.
func SetOutput(w io.Writer) Control {
	return logger.SetOutput(w)
}

func WithTimestamp() Control {
	return logger.WithTimestamp()
}

func WithCaller() Control {
	return logger.WithCaller()
}

func WithCallerWithSkipFrameCount(skipFrameCount int) Control {
	return logger.WithCallerWithSkipFrameCount(skipFrameCount)
}

func WithStack() Control {
	return logger.WithStack()
}

func WithFields(fields ...interface{}) Control {
	return logger.WithFields(fields...)
}

// GetLogger returns the current global
func DefaultLogger() FullLogger {
	return logger
}

func Debug(v ...interface{}) {
	logger.Debug(v...)
}

func Info(v ...interface{}) {
	logger.Info(v...)
}

func Warn(v ...interface{}) {
	logger.Warn(v...)
}

func Error(v ...interface{}) {
	logger.Error(v...)
}

func Fatal(v ...interface{}) {
	logger.Fatal(v...)
}

// Debugf calls the default logger's Debugf method.
func Debugf(format string, v ...interface{}) {
	logger.Debugf(format, v...)
}

// Infof calls the default logger's Infof method.
func Infof(format string, v ...interface{}) {
	logger.Infof(format, v...)
}

// Warnf calls the default logger's Warnf method.
func Warnf(format string, v ...interface{}) {
	logger.Warnf(format, v...)
}

// Errorf calls the default logger's Errorf method.
func Errorf(format string, v ...interface{}) {
	logger.Errorf(format, v...)
}

// Fatalf calls the default logger's Fatalf method and then os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	logger.Fatalf(format, v...)
}

var _ FullLogger = (*defaultLogger)(nil)

const defaultLoggerCallerSkipFrameCount = 3

type defaultLogger struct {
	log *zerolog
}

func (ll *defaultLogger) Clone() FullLogger {
	l2 := new(defaultLogger)
	l2.log = ll.log.Clone()
	return l2
}

func (ll *defaultLogger) SetOutput(w io.Writer) Control {
	ll.log.SetOutput(w)
	return ll
}

func (ll *defaultLogger) SetLevel(lv Level) Control {
	ll.log.SetLevel(lv)
	return ll
}

func (ll *defaultLogger) WithTimestamp() FullLogger {
	ll.log.WithTimestamp()
	return ll
}

func (ll *defaultLogger) WithCaller() FullLogger {
	ll.log.WithCallerWithSkipFrameCount(defaultLoggerCallerSkipFrameCount)
	return ll
}

func (ll *defaultLogger) WithCallerWithSkipFrameCount(skipFrameCount int) FullLogger {
	ll.log.WithCallerWithSkipFrameCount(defaultLoggerCallerSkipFrameCount + skipFrameCount)
	return ll
}

func (ll *defaultLogger) WithStack() FullLogger {
	ll.log.WithStack()
	return ll
}

func (ll *defaultLogger) WithFields(fields ...interface{}) FullLogger {
	ll.log.WithFields(fields...)
	return ll
}

func (ll *defaultLogger) Log(lv Level, v ...interface{}) error {
	ll.log.Log(lv, v...)
	return nil
}

func (ll *defaultLogger) Debug(v ...interface{}) {
	ll.Log(LevelDebug, v...)
}

func (ll *defaultLogger) Info(v ...interface{}) {
	ll.Log(LevelInfo, v...)
}

func (ll *defaultLogger) Warn(v ...interface{}) {
	ll.Log(LevelWarn, v...)
}

func (ll *defaultLogger) Error(v ...interface{}) {
	ll.Log(LevelError, v...)
}

func (ll *defaultLogger) Fatal(v ...interface{}) {
	ll.Log(LevelFatal, v...)
}

func (ll *defaultLogger) Debugf(format string, v ...interface{}) {
	ll.Log(LevelDebug, fmt.Sprintf(format, v...))
}

func (ll *defaultLogger) Infof(format string, v ...interface{}) {
	ll.Log(LevelInfo, fmt.Sprintf(format, v...))
}

func (ll *defaultLogger) Warnf(format string, v ...interface{}) {
	ll.Log(LevelWarn, fmt.Sprintf(format, v...))
}

func (ll *defaultLogger) Errorf(format string, v ...interface{}) {
	ll.Log(LevelError, fmt.Sprintf(format, v...))
}

func (ll *defaultLogger) Fatalf(format string, v ...interface{}) {
	ll.Log(LevelFatal, fmt.Sprintf(format, v...))
}

package slog

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test package level functions without format
func normalOutput(t *testing.T, testLevel Level, want string, args ...interface{}) {
	buf := new(bytes.Buffer)
	SetOutput(buf)
	defer SetOutput(os.Stderr)
	switch testLevel {
	case LevelDebug:
		Debug(args...)
		assert.Contains(t, buf.String(), want)
	case LevelInfo:
		Info(args...)
		assert.Contains(t, buf.String(), want)
	case LevelWarn:
		Warn(args...)
		assert.Contains(t, buf.String(), want)
	case LevelError:
		Error(args...)
		assert.Contains(t, buf.String(), want)
	case LevelFatal:
		t.Fatal("fatal method cannot be tested")
	default:
		t.Errorf("unknow level: %d", testLevel)
	}
}

// test package level functions with 'format'
func formatOutput(t *testing.T, testLevel Level, want, format string, args ...interface{}) {
	buf := new(bytes.Buffer)
	SetOutput(buf)
	defer SetOutput(os.Stderr)
	switch testLevel {
	case LevelDebug:
		Debugf(format, args...)
		assert.Contains(t, buf.String(), want)
	case LevelInfo:
		Infof(format, args...)
		assert.Contains(t, buf.String(), want)
	case LevelWarn:
		Warnf(format, args...)
		assert.Contains(t, buf.String(), want)
	case LevelError:
		Errorf(format, args...)
		assert.Contains(t, buf.String(), want)
	case LevelFatal:
		t.Fatal("fatal method cannot be tested")
	default:
		t.Errorf("unknow level: %d", testLevel)
	}
}

func TestOutput(t *testing.T) {
	defer SetLevel(LevelInfo)
	tests := []struct {
		format      string
		args        []interface{}
		testLevel   Level
		loggerLevel Level
		want        string
	}{
		{"%s %s", []interface{}{"LevelInfo", "test"}, LevelInfo, LevelWarn, ""},
		{"%s%s", []interface{}{"LevelDebug", "Test"}, LevelDebug, LevelDebug, "Test"},
		{"%s", []interface{}{"LevelError test"}, LevelError, LevelInfo, "LevelError test"},
		{"%s", []interface{}{"LevelWarn test"}, LevelWarn, LevelWarn, "LevelWarn test"},
	}

	for _, tt := range tests {
		SetLevel(tt.loggerLevel)
		normalOutput(t, tt.testLevel, tt.want, tt.args...)
		formatOutput(t, tt.testLevel, tt.want, tt.format, tt.args...)
	}
}

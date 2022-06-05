package slog

import (
	"context"
	"testing"
)

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want Level
	}{
		{
			name: "DEBUG",
			want: LevelDebug,
			s:    "DEBUG",
		},
		{
			name: "INFO",
			want: LevelInfo,
			s:    "INFO",
		},
		{
			name: "WARN",
			want: LevelWarn,
			s:    "WARN",
		},
		{
			name: "ERROR",
			want: LevelError,
			s:    "ERROR",
		},
		{
			name: "FATAL",
			want: LevelFatal,
			s:    "FATAL",
		},
		{
			name: "other",
			want: LevelInfo,
			s:    "other",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := ParseLevel(tt.s); got != tt.want {
				t.Errorf("ParseLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContext(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger1 := FromContext(ctx)
	if logger1 == nil {
		t.Fatal("expected logger to never be nil")
	}

	ctx = WithLogger(ctx, logger1)

	logger2 := FromContext(ctx)
	if logger1 != logger2 {
		t.Errorf("expected %#v to be %#v", logger1, logger2)
	}
}

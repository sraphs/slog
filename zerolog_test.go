package slog

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_zerolog_Log(t *testing.T) {
	TimestampFunc = func() time.Time {
		return time.Date(2001, time.February, 3, 4, 5, 6, 7, time.UTC)
	}
	tests := []struct {
		name string
		lv   Level
		kvs  []interface{}
		want string
	}{
		{
			name: "log with no kvs",
			lv:   LevelInfo,
			kvs:  nil,
			want: "",
		},
		{
			name: "log with empty kvs",
			lv:   LevelInfo,
			kvs:  []interface{}{},

			want: "",
		},
		{
			name: "log with single string",
			lv:   LevelInfo,
			kvs: []interface{}{
				"hello",
			},

			want: "{\"level\":\"info\",\"msg\":\"hello\"}\n",
		},
		{
			name: "log with pair string",
			lv:   LevelInfo,
			kvs: []interface{}{
				"foo",
				"bar",
			},

			want: "{\"level\":\"info\",\"foo\":\"bar\"}\n",
		},
		{
			name: "log with odd string",
			lv:   LevelInfo,
			kvs: []interface{}{
				"foo",
				"bar",
				"baz",
			},

			want: "{\"level\":\"info\",\"msg\":\"foo\",\"bar\":\"baz\"}\n",
		},
		{
			name: "log with complex",
			lv:   LevelInfo,
			kvs: []interface{}{
				"string",
				"num",
				1,
				"bool",
				true,
				"nil",
				nil,
			},
			want: "{\"level\":\"info\",\"msg\":\"string\",\"num\":1,\"bool\":true,\"nil\":null}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			z := newZerolog(w)
			assert.NoError(t, z.Log(tt.lv, tt.kvs...))
			assert.Contains(t, w.String(), tt.want)
		})
	}
}

func Test_zerolog_SetLevel(t *testing.T) {
	tests := []struct {
		name string
		lv   Level
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			z := newZerolog(w)
			z.SetLevel(tt.lv)
		})
	}
}

func Test_zerolog_SetOutput(t *testing.T) {
	z := newZerolog(nil)
	w := &bytes.Buffer{}
	z.SetOutput(w)
	z.Log(LevelInfo, "test")
	assert.Contains(t, w.String(), "test")
}

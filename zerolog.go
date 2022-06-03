package slog

import (
	"encoding/json"
	"io"
	"strconv"
	"sync"
	"time"

	zlog "github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func newZerolog(w io.Writer) *zerolog {
	initialize()
	return &zerolog{
		log:   zlog.New(w),
		level: LevelInfo,
		w:     w,
	}
}

var _ KLogger = (*zerolog)(nil)
var _ Control = (*zerolog)(nil)

// zerolog implements a Logger interface using zerolog.
type zerolog struct {
	log   zlog.Logger
	level Level
	mu    sync.Mutex
	w     io.Writer
}

func (z *zerolog) Log(lv Level, kvs ...interface{}) error {
	if z.level > lv {
		return nil
	}

	if len(kvs) == 0 {
		return nil
	}

	var e *zlog.Event

	switch lv {
	case LevelDebug:
		e = z.log.Debug()
	case LevelInfo:
		e = z.log.Info()
	case LevelWarn:
		e = z.log.Warn()
	case LevelError:
		e = z.log.Error()
	case LevelFatal:
		e = z.log.Fatal()
	default:
		e = z.log.Info()
	}

	for i, v := range kvs {
		if err, ok := v.(error); ok {
			e.Err(err)
			kvs = append(kvs[:i], kvs[i+1:]...)
		}
	}

	if len(kvs) == 0 {
		e.Send()
		return nil
	}

	if len(kvs)%2 == 0 {
		e.Fields(kvs)
	} else {
		if s, ok := kvs[0].(string); ok {
			e.Str(MessageFieldName, s)
		} else {
			e.Fields([]interface{}{MessageFieldName, kvs[0]})
		}

		if len(kvs) > 1 {
			e.Fields(kvs[1:])
		}
	}

	e.Send()

	return nil
}

// SetLevel sets the current global log level.
func (z *zerolog) SetLevel(l Level) Control {
	z.mu.Lock()
	defer z.mu.Unlock()
	z.level = l
	return z
}

func (z *zerolog) SetOutput(w io.Writer) Control {
	z.mu.Lock()
	defer z.mu.Unlock()
	z.log = z.log.Output(w)
	return z
}

func (z *zerolog) Clone() *zerolog {
	z2 := newZerolog(z.w)
	z2.mu.Lock()
	defer z2.mu.Unlock()
	z2.level = z.level
	z2.log = z.log.With().Logger()
	z2.w = z.w
	return z2
}

func (z *zerolog) WithTimestamp() *zerolog {
	z.mu.Lock()
	defer z.mu.Unlock()
	z.log = z.log.With().Timestamp().Logger()
	return z
}

func (z *zerolog) WithCaller() *zerolog {
	z.mu.Lock()
	defer z.mu.Unlock()
	z.log = z.log.With().Caller().Logger()
	return z
}

func (z *zerolog) WithCallerWithSkipFrameCount(skipFrameCount int) *zerolog {
	z.mu.Lock()
	defer z.mu.Unlock()
	z.log = z.log.With().CallerWithSkipFrameCount(CallerSkipFrameCount + skipFrameCount).Logger()
	return z
}

func (z *zerolog) WithStack() *zerolog {
	z.mu.Lock()
	defer z.mu.Unlock()
	z.log = z.log.With().Stack().Logger()
	return z
}

func (z *zerolog) WithFields(fields ...interface{}) *zerolog {
	z.mu.Lock()
	defer z.mu.Unlock()
	z.log = z.log.With().Fields(fields).Logger()
	return z
}

var (
	MultiLevelWriter = zlog.MultiLevelWriter
)

var (
	// TimestampFieldName is the field name used for the timestamp field.
	TimestampFieldName = "ts"

	// LevelFieldName is the field name used for the level field.
	LevelFieldName = "level"

	// LevelTraceValue is the value used for the trace level field.
	LevelTraceValue = "trace"
	// LevelDebugValue is the value used for the debug level field.
	LevelDebugValue = "debug"
	// LevelInfoValue is the value used for the info level field.
	LevelInfoValue = "info"
	// LevelWarnValue is the value used for the warn level field.
	LevelWarnValue = "warn"
	// LevelErrorValue is the value used for the error level field.
	LevelErrorValue = "error"
	// LevelFatalValue is the value used for the fatal level field.
	LevelFatalValue = "fatal"
	// LevelPanicValue is the value used for the panic level field.
	LevelPanicValue = "panic"

	// LevelFieldMarshalFunc allows customization of global level field marshaling.
	LevelFieldMarshalFunc = func(l zlog.Level) string {
		return l.String()
	}

	// MessageFieldName is the field name used for the message field.
	MessageFieldName = "msg"

	// ErrorFieldName is the field name used for error fields.
	ErrorFieldName = "error"

	// CallerFieldName is the field name used for caller field.
	CallerFieldName = "caller"

	// CallerSkipFrameCount is the number of stack frames to skip to find the caller.
	CallerSkipFrameCount = 2 + 1

	// CallerMarshalFunc allows customization of global caller marshaling
	CallerMarshalFunc = func(file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}

	// ErrorStackFieldName is the field name used for error stacks.
	ErrorStackFieldName = "stack"

	// ErrorStackMarshaler extract the stack from err if any.
	ErrorStackMarshaler = pkgerrors.MarshalStack

	// ErrorMarshalFunc allows customization of global error marshaling
	ErrorMarshalFunc = func(err error) interface{} {
		return err
	}

	// InterfaceMarshalFunc allows customization of interface marshaling.
	// Default: "encoding/json.Marshal"
	InterfaceMarshalFunc = json.Marshal

	// TimeFieldFormat defines the time format of the Time field type. If set to
	// TimeFormatUnix, TimeFormatUnixMs or TimeFormatUnixMicro, the time is formatted as an UNIX
	// timestamp as integer.
	TimeFieldFormat = time.RFC3339

	// TimestampFunc defines the function called to generate a timestamp.
	TimestampFunc = time.Now

	// DurationFieldUnit defines the unit for time.Duration type fields added
	// using the Dur method.
	DurationFieldUnit = time.Millisecond

	// DurationFieldInteger renders Dur fields as integer instead of float if
	// set to true.
	DurationFieldInteger = false

	// ErrorHandler is called whenever zerolog fails to write an event on its
	// output. If not set, an error is printed on the stderr. This handler must
	// be thread safe and non-blocking.
	ErrorHandler func(err error)
)

func initialize() {
	zlog.TimestampFieldName = TimestampFieldName
	zlog.LevelFieldName = LevelFieldName
	zlog.LevelTraceValue = LevelTraceValue
	zlog.LevelDebugValue = LevelDebugValue
	zlog.LevelInfoValue = LevelInfoValue
	zlog.LevelWarnValue = LevelWarnValue
	zlog.LevelErrorValue = LevelErrorValue
	zlog.LevelFatalValue = LevelFatalValue
	zlog.LevelPanicValue = LevelPanicValue
	zlog.LevelFieldMarshalFunc = LevelFieldMarshalFunc
	zlog.MessageFieldName = MessageFieldName
	zlog.ErrorFieldName = ErrorFieldName
	zlog.CallerFieldName = CallerFieldName
	zlog.CallerSkipFrameCount = CallerSkipFrameCount
	zlog.CallerMarshalFunc = CallerMarshalFunc
	zlog.ErrorStackFieldName = ErrorStackFieldName
	zlog.ErrorStackMarshaler = ErrorStackMarshaler
	zlog.ErrorMarshalFunc = ErrorMarshalFunc
	zlog.InterfaceMarshalFunc = InterfaceMarshalFunc
	zlog.TimeFieldFormat = TimeFieldFormat
	zlog.TimestampFunc = TimestampFunc
	zlog.DurationFieldUnit = DurationFieldUnit
	zlog.DurationFieldInteger = DurationFieldInteger
	zlog.ErrorHandler = ErrorHandler
}

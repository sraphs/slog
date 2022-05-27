package slog_test

import (
	"time"

	"github.com/pkg/errors"

	"github.com/sraphs/slog"
)

func Example() {
	slog.TimestampFunc = func() time.Time {
		return time.Date(2001, time.February, 3, 4, 5, 6, 7, time.UTC)
	}

	slog.Init(&slog.Config{
		Level:   "info",
		Path:    "log/app.log",
		MaxSize: 100,
		MaxAge:  7,
	}).WithTimestamp().WithCaller().WithStack()

	// level log
	slog.Debug("hello world")
	slog.Info("hello world")
	slog.Warn("hello world")
	slog.Error("hello world")
	// slog.Fatal("hello world")

	// format log
	slog.Debugf("hello %s", "world")
	slog.Infof("hello %s", "world")
	slog.Warnf("hello %s", "world")
	slog.Errorf("hello %s", "world")
	// slog.Fatalf("hello %s", "world")

	// log err with stack
	err := outer()
	slog.Error(err)

	slog.Clone().WithFields("foo", "bar").Info("hello world")

	slog.Info("hello world")

	// Outputs:
	// {"level":"info","msg":"hello world","ts":"2001-02-03T04:05:06Z","caller":"example_test.go:25"}
	// {"level":"warn","msg":"hello world","ts":"2001-02-03T04:05:06Z","caller":"example_test.go:26"}
	// {"level":"error","msg":"hello world","ts":"2001-02-03T04:05:06Z","caller":"example_test.go:27"}
	// {"level":"info","msg":"hello world","ts":"2001-02-03T04:05:06Z","caller":"example_test.go:32"}
	// {"level":"warn","msg":"hello world","ts":"2001-02-03T04:05:06Z","caller":"example_test.go:33"}
	// {"level":"error","msg":"hello world","ts":"2001-02-03T04:05:06Z","caller":"example_test.go:34"}
	// {"level":"error","stack":[{"func":"inner","line":"58","source":"example_test.go"},{"func":"middle","line":"62","source":"example_test.go"},{"func":"outer","line":"70","source":"example_test.go"},{"func":"Example","line":"38","source":"example_test.go"},{"func":"runExample","line":"63","source":"run_example.go"},{"func":"runExamples","line":"44","source":"example.go"},{"func":"(*M).Run","line":"1721","source":"testing.go"},{"func":"main","line":"61","source":"_testmain.go"},{"func":"main","line":"250","source":"proc.go"},{"func":"goexit","line":"1571","source":"asm_amd64.s"}],"error":"seems we have an error here","ts":"2001-02-03T04:05:06Z","caller":"example_test.go:39"}
	// {"level":"info","foo":"bar","msg":"hello world","ts":"2001-02-03T04:05:06Z","caller":"run_example.go:63"}
	// {"level":"info","msg":"hello world","ts":"2001-02-03T04:05:06Z","caller":"example_test.go:43"}
}

func inner() error {
	return errors.New("seems we have an error here")
}

func middle() error {
	err := inner()
	if err != nil {
		return err
	}
	return nil
}

func outer() error {
	err := middle()
	if err != nil {
		return err
	}
	return nil
}

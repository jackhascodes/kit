package log_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/jackhascodes/kit/log"
)

func simulatedExit(code int) {
	fmt.Println()
	fmt.Println("Process finished with exit code ", code)
}

func ExampleWithDefaultLogLevel() {
	info := log.InitLog()
	trace := log.InitLog(log.WithDefaultLogLevel(log.Trace))
	info.Print("foo")
	trace.Print("bar")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = info; message = foo
	// time = 2020-06-01T00:00:00Z00:00; level = trace; message = bar
}

func ExampleWithExit() {
	logger := log.InitLog(
		log.WithExit(func(code int) {
			// slack.Post(fmt.Sprintf("example application died with code %d", code)
			fmt.Println("notified slack")
			simulatedExit(code)
		}))
	logger.Fatal(2, "foo")

	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = fatal; message = foo; stack = [example.go@main:10 example.go@main:4]
	// notified slack
	//
	// Process finished with exit code  2
}

func ExampleWithFieldOrder() {
	logger := log.InitLog(
		log.WithFields(
			log.KV{"foo", "bar"},
			log.KV{"baz", "qux"},
		),
		log.WithFieldOrder("baz", "time", "message"),
	)
	logger.Debug("quux")
	// Output:
	// baz = qux; time = 2020-06-01T00:00:00Z00:00; message = quux; level = debug; foo = bar
}

func ExampleWithFields() {
	logger := log.InitLog(
		log.WithFields(
			log.KV{"foo", "bar"},
			log.KV{"baz", "qux"},
		))
	logger.Debug("quux")
	logger.Info("gurp")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = debug; message = quux; foo = bar; baz = qux
	// time = 2020-06-01T00:00:00Z00:00; level = info; message = gurp; foo = bar; baz = qux

}

func ExampleWithFormatter() {
	logger := log.InitLog(log.WithFormatter(&log.JsonFormatter{}))
	logger.Debug("foo")
	// Output:
	// {"time":"2020-06-01T00:00:00Z00:00","level":"debug","message":"foo"}
}

func ExampleWithLevelNames() {
	names := map[log.LogLevel]string{
		log.Warn:  "*WARN*",
		log.Fatal: "_*FATAL*_",
		log.Panic: "!_*PANIC*_!",
	}
	logger := log.InitLog(log.WithLevelNames(names))
	logger.Debug("foo")
	logger.Warn("bar")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = debug; message = foo
	// time = 2020-06-01T00:00:00Z00:00; level = *WARN*; message = bar
}

func ExampleWithMinLevel() {
	logger := log.InitLog(log.WithMinLevel(log.Warn))
	logger.Debug("foo")
	logger.Info("bar")
	logger.Warn("baz")
	// Output:
	// 	time = 2020-06-01T00:00:00Z00:00; level = warn; message = baz
}

func ExampleWithOutput() {
	var buf bytes.Buffer
	logger := log.InitLog(log.WithOutput(&buf))
	logger.Debug("foo")
	// Output:
}

func ExampleWithPrefix() {
	logger := log.InitLog(log.WithPrefix("DEV: "))
	logger.Debug("foo")
	logger.Info("bar")
	// Output:
	// DEV: time = 2020-06-01T00:00:00Z00:00; level = debug; message = foo
	// DEV: time = 2020-06-01T00:00:00Z00:00; level = info; message = bar
}

func ExampleWithPreProcessors() {
	logger := log.InitLog(log.WithPreProcessors(func(f *log.Fields) {
		if err, _ := f.Get("error"); err != nil {
			// errorStore.save(err)
			fmt.Println(fmt.Sprintf("saved error '%s' from fields", err))
		}
	}))
	logger.Debug("foo")
	err := errors.New("bar")
	logger.Error(err, "baz")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = debug; message = foo
	// saved error 'bar' from fields
	// time = 2020-06-01T00:00:00Z00:00; level = error; message = baz; error = bar; stack = [example.go@main:10 example.go@main:4]
}

func ExampleWithPostProcessors() {
	logger := log.InitLog(log.WithPostProcessors(func(f *log.Fields) {
		if _, err := f.Get("error"); err != nil {
			// errorStore.save(err)
			fmt.Println(fmt.Sprintf("saved error: %s from fields: %v", err, f))
		}
	}))
	logger.Debug("foo")
	err := errors.New("bar")
	logger.Error(err, "baz")
	// Outputs:
	// 	time = 2020-06-01T00:00:00Z00:00; level = debug; message = foo
	// 	saved error fields {error: bar, time: 2020-06-01T00:00:00Z00:00, level: error; message: baz}
	// 	time = 2020-06-01T00:00:00Z00:00; level = error; message = baz; error = bar
}

func ExampleWithTimer() {
	logger := log.InitLog(
		log.WithTimer(func(format string) string {
			// return time.Now().Format(format)
			return "2121-06-01T00:00:00Z00:00"
		}))
	logger.Debug("foo")
	// Output:
	// 	time = 2121-06-01T00:00:00Z00:00; level = debug; message = foo
}

func ExampleWithTimer_countdown() {
	targetTime := time.Now().Add(10 * time.Second)
	logger := log.InitLog(log.WithTimer(
		func(format string) string {
			remaining := targetTime.Sub(time.Now())
			return fmt.Sprintf(format, int(math.Round(remaining.Seconds())))
		}),
		log.WithTimeFormat("%d seconds remain"),
	)
	logger.Debug("foo")
	time.Sleep(1 * time.Second)
	logger.Debug("bar")
	// Output:
	// time = 10 seconds remain; level = debug; message = foo
	// time = 9 seconds remain; level = debug; message = bar
}

func ExampleWithTimeFormat() {
	logger := log.InitLog(log.WithTimeFormat(time.RFC822))
	logger.Debug("foo")
	// Outputs:
	// time = 02 Jun 20 00:00 UTC; level = debug; message = foo

}

func TestWithFormatter(t *testing.T) {
	logger := log.InitLog(
		log.WithFormatter(&log.JsonFormatter{}),
		log.WithExit(func(c int) {}),
	)
	tests := []struct {
		name  string
		doLog func(*log.Log)
		want  string
	}{
		{"print",
			func(l *log.Log) { l.Print("logger output") },
			`{"time":"2020-06-01T00:00:00Z00:00","level":"info","message":"logger output"}` + "\n",
		}, {"trace",
			func(l *log.Log) { l.Trace("logger output") },
			`{"time":"2020-06-01T00:00:00Z00:00","level":"trace","message":"logger output"}` + "\n",
		}, {"debug",
			func(l *log.Log) { l.Debug("logger output") },
			`{"time":"2020-06-01T00:00:00Z00:00","level":"debug","message":"logger output"}` + "\n",
		}, {"info",
			func(l *log.Log) { l.Info("logger output") },
			`{"time":"2020-06-01T00:00:00Z00:00","level":"info","message":"logger output"}` + "\n",
		}, {"warn",
			func(l *log.Log) { l.Warn("logger output") },
			`{"time":"2020-06-01T00:00:00Z00:00","level":"warn","message":"logger output"}` + "\n",
		}, {"error",
			func(l *log.Log) { l.Error(errors.New("oops"), "logger output") },
			`{"time":"2020-06-01T00:00:00Z00:00","level":"error","message":"logger output","error":"oops","stack":["example.go@main:10","example.go@main:4"]}` + "\n",
		}, {"printf",
			func(l *log.Log) { l.Printf("logger output %d %s", 2, "with vars") },
			`{"time":"2020-06-01T00:00:00Z00:00","level":"info","message":"logger output 2 with vars"}` + "\n",
		}, {"tracef",
			func(l *log.Log) { l.Tracef("logger output %d %s", 2, "with vars") },
			`{"time":"2020-06-01T00:00:00Z00:00","level":"trace","message":"logger output 2 with vars"}` + "\n",
		}, {"debugf",
			func(l *log.Log) { l.Debugf("logger output %d %s", 2, "with vars") },
			`{"time":"2020-06-01T00:00:00Z00:00","level":"debug","message":"logger output 2 with vars"}` + "\n",
		}, {"infof",
			func(l *log.Log) { l.Infof("logger output %d %s", 2, "with vars") },
			`{"time":"2020-06-01T00:00:00Z00:00","level":"info","message":"logger output 2 with vars"}` + "\n",
		}, {"warnf",
			func(l *log.Log) { l.Warnf("logger output %d %s", 2, "with vars") },
			`{"time":"2020-06-01T00:00:00Z00:00","level":"warn","message":"logger output 2 with vars"}` + "\n",
		}, {"errorf",
			func(l *log.Log) { l.Errorf(errors.New("oops"), "logger output %d %s", 2, "with vars") },

			`{"time":"2020-06-01T00:00:00Z00:00","level":"error","message":"logger output 2 with vars","error":"oops","stack":["example.go@main:10","example.go@main:4"]}` + "\n",
		}, {"fatal",
			func(l *log.Log) { l.Fatal(1, "logger output") },
			`{"time":"2020-06-01T00:00:00Z00:00","level":"fatal","message":"logger output","stack":["example.go@main:10","example.go@main:4"]}` + "\n",
		},
		{
			"fatalf",
			func(l *log.Log) { l.Fatalf(1, "logger output %d %s", 2, "with vars") },
			`{"time":"2020-06-01T00:00:00Z00:00","level":"fatal","message":"logger output 2 with vars","stack":["example.go@main:10","example.go@main:4"]}` + "\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger.SetOutput(&buf)
			tt.doLog(logger)
			got := buf.String()
			if got != tt.want {
				t.Errorf("wanted:\n\t'%s'\ngot:\n\t'%s'", tt.want, got)
			}
		})
	}
}

func TestWithFields(t *testing.T) {
	logger := log.InitLog(
		log.WithFormatter(&log.JsonFormatter{}),
		log.WithFields(log.KV{"testing", "on"}, log.KV{"second", 100}),
	)
	want := `{"time":"2020-06-01T00:00:00Z00:00","level":"debug","message":"logger output","testing":"on","second":100}` + "\n"
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	logger.Debug("logger output")
	got := buf.String()
	if got != want {
		t.Errorf("expected:\n\t'%s'\nto be:\n\t'%s'", got, want)
	}
}

func TestWithFieldOrder(t *testing.T) {
	logger := log.InitLog(
		log.WithFormatter(&log.JsonFormatter{}),
		log.WithFieldOrder("level", "message", "time"),
	)
	want := `{"level":"debug","message":"logger output","time":"2020-06-01T00:00:00Z00:00"}` + "\n"
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	logger.Debug("logger output")
	got := buf.String()
	if got != want {
		t.Errorf("expected:\n\t'%s'\nto be:\n\t'%s'", got, want)
	}
}

func TestWithPrePostProcessors(t *testing.T) {
	var (
		primary bytes.Buffer
		pre     bytes.Buffer
		post    bytes.Buffer
	)

	logger := log.InitLog(
		log.WithFormatter(&log.JsonFormatter{}),
		log.WithPreProcessors(func(f *log.Fields) {
			b, _ := json.Marshal(f)
			pre.Write(b)
		}),
		log.WithPostProcessors(func(f *log.Fields) {
			// write a bytestream logger.Without appending a newline char.
			b, _ := json.Marshal(f)
			post.Write(b)
		}),
		log.WithFields(log.KV{"test", "case"}),
		log.WithOutput(&primary),
	)
	logger.Debug()
	primaryOut := primary.String()
	preOut := pre.String()
	postOut := post.String()
	primary.Truncate(0)
	pre.Truncate(0)
	post.Truncate(0)
	preWant := `{"test":"case"}`
	primaryWant := `{"time":"2020-06-01T00:00:00Z00:00","level":"debug","message":"","test":"case"}` + "\n"
	postWant := `{"time":"2020-06-01T00:00:00Z00:00","level":"debug","message":"","test":"case"}`
	if primaryOut != primaryWant {
		t.Errorf("primary expected:\n\t'%s' to be:\n\t'%s'", primaryOut, primaryWant)
	}
	if preOut != preWant {
		t.Errorf("pre expected:\n\t'%s' to be:\n\t'%s'", preOut, preWant)
	}
	if postOut != postWant {
		t.Errorf("post expected:\n\t'%s' to be:\n\t'%s'", postOut, postWant)
	}
}

func TestWithLevelNames(t *testing.T) {
	var buf bytes.Buffer
	logger := log.InitLog(
		log.WithFormatter(&log.JsonFormatter{}),
		log.WithLevelNames(map[log.LogLevel]string{
			log.Debug: "DEBUG",
			log.Info:  "INFO",
			log.Warn:  "WARN",
		}),
		log.WithOutput(&buf),
	)
	logger.Debug()
	logger.Info()
	logger.Warn()
	want := `{"time":"2020-06-01T00:00:00Z00:00","level":"DEBUG","message":""}` + "\n" +
		`{"time":"2020-06-01T00:00:00Z00:00","level":"INFO","message":""}` + "\n" +
		`{"time":"2020-06-01T00:00:00Z00:00","level":"WARN","message":""}` + "\n"
	got := buf.String()
	if got != want {
		t.Errorf("expected:\n\t'%s' to be:\n\t'%s'", got, want)
	}

}

func TestWithMinLevel(t *testing.T) {
	tests := []struct {
		name     string
		minLevel log.LogLevel
		want     string
	}{
		{"trace", log.Trace,
			`{"time":"2020-06-01T00:00:00Z00:00","level":"trace","message":""}` + "\n" +
				`{"time":"2020-06-01T00:00:00Z00:00","level":"debug","message":""}` + "\n" +
				`{"time":"2020-06-01T00:00:00Z00:00","level":"info","message":""}` + "\n" +
				`{"time":"2020-06-01T00:00:00Z00:00","level":"warn","message":""}` + "\n",
		}, {"debug", log.Debug,
			`{"time":"2020-06-01T00:00:00Z00:00","level":"debug","message":""}` + "\n" +
				`{"time":"2020-06-01T00:00:00Z00:00","level":"info","message":""}` + "\n" +
				`{"time":"2020-06-01T00:00:00Z00:00","level":"warn","message":""}` + "\n",
		}, {"info", log.Info,
			`{"time":"2020-06-01T00:00:00Z00:00","level":"info","message":""}` + "\n" +
				`{"time":"2020-06-01T00:00:00Z00:00","level":"warn","message":""}` + "\n",
		}, {"warn", log.Warn,
			`{"time":"2020-06-01T00:00:00Z00:00","level":"warn","message":""}` + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := log.InitLog(
				log.WithFormatter(&log.JsonFormatter{}),
				log.WithMinLevel(tt.minLevel),
				log.WithOutput(&buf),
			)
			logger.Trace()
			logger.Debug()
			logger.Info()
			logger.Warn()
			got := buf.String()
			buf.Truncate(0)
			if got != tt.want {
				t.Errorf("expected:\n\t'%s'\nto be:\n\t'%s'", got, tt.want)
			}
		})
	}
}

func TestWithTimer_countdown(t *testing.T) {
	var buf bytes.Buffer
	countdownTo := testNow.Add(10 * time.Second)
	logger := log.InitLog(
		log.WithTimer(func(string) string {
			remaining := countdownTo.Sub(testNow)
			return fmt.Sprintf("%d seconds remaining", int(remaining.Seconds()))
		}),
		log.WithFormatter(&log.JsonFormatter{}),
		log.WithFieldOrder("time", "level", "message"),
		log.WithOutput(&buf),
	)
	logger.Debug()
	want := `{"time":"10 seconds remaining","level":"debug","message":""}` + "\n"
	got := buf.String()
	if got != want {
		t.Errorf("expected:\n\t'%s'\nto be:\n\t'%s'", got, want)
	}
}

func TestWithPrefix(t *testing.T) {
	var buf bytes.Buffer
	logger := log.InitLog(
		log.WithFormatter(&log.JsonFormatter{}),
		log.WithFieldOrder("time", "level", "message"),
		log.WithPrefix("WITH_PREFIX: "),
		log.WithOutput(&buf),
	)
	want := `WITH_PREFIX: {"time":"2020-06-01T00:00:00Z00:00","level":"debug","message":""}` + "\n" +
		`WITH_PREFIX: {"time":"2020-06-01T00:00:00Z00:00","level":"debug","message":""}` + "\n"
	logger.Debug()
	logger.Debug()
	got := buf.String()
	if got != want {
		t.Errorf("expected:\n\t'%s'\nto be:\n\t'%s'", got, want)
	}
}

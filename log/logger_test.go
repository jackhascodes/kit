package log_test

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jackhascodes/kit/log"
)

var (
	defaultTimer = log.DefaultTimer(func(format string) string { return "2020-06-01T00:00:00Z00:00" })
	defaultStack = log.DefaultStackMethod(func() []string { return []string{"example.go@main:10", "example.go@main:4"} })
	defaultExit  = log.DefaultExit(func(code int) { fmt.Println(fmt.Sprintf("\nProcess finished with exit code %d", code)) })
	defaultPanic = log.DefaultPanic(func(v interface{}) {
		fmt.Println("panic:", v)
		fmt.Println("\ngoroutine 1 [running]:")
		fmt.Println("main.main()")
		fmt.Println("	/some/project/main.go:11 +0x186")
	})
)

func ExampleLog() {
	logger := log.InitLog()
	logger.Trace("foo")
	logger.Debug("bar")
	logger.Info("baz")
	logger.Warn("qux")
	logger.Error(errors.New("quux"), "quuz")
	logger.Fatal(2, "corge")
	logger.Panic("grault")

	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = trace; message = foo
	// time = 2020-06-01T00:00:00Z00:00; level = debug; message = bar
	// time = 2020-06-01T00:00:00Z00:00; level = info; message = baz
	// time = 2020-06-01T00:00:00Z00:00; level = warn; message = qux
	// time = 2020-06-01T00:00:00Z00:00; level = error; message = quuz; error = quux; stack = [example.go@main:10 example.go@main:4]
	// time = 2020-06-01T00:00:00Z00:00; level = fatal; message = corge; stack = [example.go@main:10 example.go@main:4]
	//
	// Process finished with exit code 2
	// time = 2020-06-01T00:00:00Z00:00; level = panic; message = grault; stack = [example.go@main:10 example.go@main:4]
	// panic: grault
	//
	// goroutine 1 [running]:
	// main.main()
	// 	/some/project/main.go:11 +0x186

}

func ExampleInitLog() {
	logger := log.InitLog()
	logger.Debug("hello world")
	logger.Info("foo")
	logger.Error(errors.New("boo"), "bar")
	logger.Fatal(2, "baz")

	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = debug; message = hello world
	// time = 2020-06-01T00:00:00Z00:00; level = info; message = foo
	// time = 2020-06-01T00:00:00Z00:00; level = error; message = bar; error = boo; stack = [example.go@main:10 example.go@main:4]
	// time = 2020-06-01T00:00:00Z00:00; level = fatal; message = baz; stack = [example.go@main:10 example.go@main:4]
	//
	// Process finished with exit code 2
}

func ExampleInitLog_withExit() {
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

func ExampleInitLog_withFieldOrder() {
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

func ExampleInitLog_withFields() {
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

func ExampleInitLog_withFormatter() {
	logger := log.InitLog(log.WithFormatter(&log.JsonFormatter{}))
	logger.Debug("foo")
	// Output:
	// {"time":"2020-06-01T00:00:00Z00:00","level":"debug","message":"foo"}
}

func ExampleInitLog_withLevelNames() {
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

func ExampleInitLog_withMinLevel() {
	logger := log.InitLog(log.WithMinLevel(log.Warn))
	logger.Debug("foo")
	logger.Info("bar")
	logger.Warn("baz")
	// Output:
	// 	time = 2020-06-01T00:00:00Z00:00; level = warn; message = baz
}

func ExampleInitLog_withOutput() {
	var buf bytes.Buffer
	logger := log.InitLog(log.WithOutput(&buf))
	logger.Debug("foo")
	// Output:
}

func ExampleInitLog_withPrefix() {
	logger := log.InitLog(log.WithPrefix("DEV: "))
	logger.Debug("foo")
	logger.Info("bar")
	// Output:
	// DEV: time = 2020-06-01T00:00:00Z00:00; level = debug; message = foo
	// DEV: time = 2020-06-01T00:00:00Z00:00; level = info; message = bar
}

func ExampleInitLog_withPreProcessors() {
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

func ExampleInitLog_withPostProcessors() {
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

func ExampleInitLog_withTimer() {
	logger := log.InitLog(
		log.WithTimer(func(format string) string {
			// return time.Now().Format(format)
			return "2121-06-01T00:00:00Z00:00"
		}))
	logger.Debug("foo")
	// Output:
	// 	time = 2121-06-01T00:00:00Z00:00; level = debug; message = foo
}

func ExampleInitLog_withTimer_countdown() {
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

func ExampleInitLog_withTimeFormat() {
	logger := log.InitLog(log.WithTimeFormat(time.RFC822))
	logger.Debug("foo")
	// Outputs:
	// time = 02 Jun 20 00:00 UTC; level = debug; message = foo

}

func ExampleLog_Debug() {
	logger := log.InitLog()
	logger.Debug("foo", "bar")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = debug; message = foobar
}

func ExampleLog_Debugf() {
	logger := log.InitLog()
	logger.Debugf("foo %s %s", "bar", "baz")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = debug; message = foo bar baz
}

func ExampleLog_Error() {
	logger := log.InitLog()
	logger.Error(errors.New("oops"), "foo")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = error; message = foo; error = oops; stack = [example.go@main:10 example.go@main:4]
}

func ExampleLog_Errorf() {
	logger := log.InitLog()
	logger.Errorf(errors.New("oops"), "foo %s %s", "bar", "baz")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = error; message = foo bar baz; error = oops; stack = [example.go@main:10 example.go@main:4]

}

func ExampleLog_Fatal() {
	logger := log.InitLog()
	logger.Fatal(2, "foo", "bar")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = fatal; message = foobar; stack = [example.go@main:10 example.go@main:4]
	//
	// Process finished with exit code 2
}

func ExampleLog_Fatalf() {
	logger := log.InitLog()
	logger.Fatalf(2, "foo %s %s", "bar", "baz")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = fatal; message = foo bar baz; stack = [example.go@main:10 example.go@main:4]
	//
	// Process finished with exit code 2
}

func ExampleLog_Info() {
	logger := log.InitLog()
	logger.Info("foo", "bar")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = info; message = foobar
}

func ExampleLog_Infof() {
	logger := log.InitLog()
	logger.Infof("foo %s %s", "bar", "baz")
	// Output:
	// 	time = 2020-06-01T00:00:00Z00:00; level = info; message = foo bar baz
}

func ExampleLog_Log() {
	logger := log.InitLog()
	logger.Log(log.Info, "foo", "bar")
	logger.Log(log.Debug, "baz")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = info; message = foobar
	// time = 2020-06-01T00:00:00Z00:00; level = debug; message = baz
}

func ExampleLog_Panic() {
	logger := log.InitLog()
	logger.Panic("foo", "bar")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = panic; message = foobar; stack = [example.go@main:10 example.go@main:4]
	// panic: foobar
	//
	// goroutine 1 [running]:
	// main.main()
	// 	/some/project/main.go:11 +0x186

}

func ExampleLog_Panicf() {
	logger := log.InitLog()
	logger.Panicf("foo %s %s", "bar", "baz")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = panic; message = foo bar baz; stack = [example.go@main:10 example.go@main:4]
	// panic: foo bar baz
	//
	// goroutine 1 [running]:
	// main.main()
	// 	/some/project/main.go:11 +0x186
}

func ExampleLog_Prefix() {
	logger := log.InitLog()
	logger.Debug("no prefix")
	logger.Prefix("A: ").Debug("foo", "bar")
	logger.Debug("baz")
	prefixed := logger.Prefix("B: ")
	prefixed.Debug("foo")
	prefixed.Debug("bar")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = debug; message = no prefix
	// A: time = 2020-06-01T00:00:00Z00:00; level = debug; message = foobar
	// time = 2020-06-01T00:00:00Z00:00; level = debug; message = baz
	// B: time = 2020-06-01T00:00:00Z00:00; level = debug; message = foo
	// B: time = 2020-06-01T00:00:00Z00:00; level = debug; message = bar
}

func ExampleLog_Print() {
	logger := log.InitLog()
	logger.Print("foo", "bar")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = info; message = foobar
}

func ExampleLog_Printf() {
	logger := log.InitLog()
	logger.Printf("foo %s %s", "bar", "baz")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = info; message = foo bar baz

}

func ExampleLog_SetOutput() {
	var buf bytes.Buffer
	logger := log.InitLog()
	logger.Debug("foo")
	logger.SetOutput(&buf)
	// "bar" message does not output. Log line is written to `buf`
	logger.Debug("bar")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = debug; message = foo
}

func ExampleLog_SetPrefix() {
	logger := log.InitLog()
	logger.Debug("no prefix")
	logger.SetPrefix("A: ")
	logger.Debug("foo")
	logger.Debug("bar")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = debug; message = no prefix
	// A: time = 2020-06-01T00:00:00Z00:00; level = debug; message = foo
	// A: time = 2020-06-01T00:00:00Z00:00; level = debug; message = bar
}

func ExampleLog_SetWatch() {
	logger := log.InitLog()
	a := "foo"
	b := map[string]string{"bar": "baz"}
	logger.SetWatch("a", &a)
	logger.SetWatch("b", &b)
	logger.Debug("qux")
	a = "quux"
	logger.Debug("qux")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = debug; message = qux; watch = map[a:foo b:map[bar:baz]]
	// time = 2020-06-01T00:00:00Z00:00; level = debug; message = qux; watch = map[a:quux b:map[bar:baz]]
}

func ExampleLog_Trace() {
	logger := log.InitLog()
	logger.Trace("foo", "bar")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = trace; message = foobar
}
func ExampleLog_Tracef() {
	logger := log.InitLog()
	logger.Tracef("foo %s %s", "bar", "baz")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = trace; message = foo bar baz

}

func ExampleLog_UnsetWatch() {
	logger := log.InitLog()
	a := "foo"
	b := "bar"
	logger.SetWatch("a", &a)
	logger.SetWatch("b", &b)
	logger.Debug("baz")
	logger.UnsetWatch("b")
	logger.Debug("qux")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = debug; message = baz; watch = map[a:foo b:bar]
	// time = 2020-06-01T00:00:00Z00:00; level = debug; message = qux; watch = map[a:foo]
}

func ExampleLog_Warn() {
	logger := log.InitLog()
	logger.Warn("foo", "bar")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = warn; message = foobar
}

func ExampleLog_Warnf() {
	logger := log.InitLog()
	logger.Warnf("foo %s %s", "bar", "baz")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = warn; message = foo bar baz
}

func ExampleLog_WithFields_inline() {
	logger := log.InitLog()
	logger.WithFields(log.InitFields(log.KV{"foo", "bar"})).Debug("baz")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = debug; message = baz; foo = bar
}

var testNow = time.Now()

func TestJsonFormatter_Format(t *testing.T) {
	f := log.InitFields().
		AddField("test", "string").
		AddField("make", []string{"slice"})
	j := []byte(`{"test":"string","make":["slice"]}`)
	tests := []struct {
		name    string
		fields  *log.Fields
		want    []byte
		wantErr bool
	}{
		{"basic", f, j, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &log.JsonFormatter{}
			got, err := j.Format(tt.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Format() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitLog(t *testing.T) {
	logger := log.InitLog()
	tests := []struct {
		name  string
		doLog func(*log.Log)
		want  string
	}{
		{"print",
			func(l *log.Log) { l.Print("log output") },
			"time = 2020-06-01T00:00:00Z00:00; level = info; message = log output\n",
		}, {"trace",
			func(l *log.Log) { l.Trace("log output") },
			"time = 2020-06-01T00:00:00Z00:00; level = trace; message = log output\n",
		}, {"debug",
			func(l *log.Log) { l.Debug("log output") },
			"time = 2020-06-01T00:00:00Z00:00; level = debug; message = log output\n",
		}, {"info",
			func(l *log.Log) { l.Info("log output") },
			"time = 2020-06-01T00:00:00Z00:00; level = info; message = log output\n",
		}, {"warn",
			func(l *log.Log) { l.Warn("log output") },
			"time = 2020-06-01T00:00:00Z00:00; level = warn; message = log output\n",
		}, {"error",
			func(l *log.Log) { l.Error(errors.New("oops"), "log output") },
			"time = 2020-06-01T00:00:00Z00:00; level = error; message = log output; error = oops; stack = [example.go@main:10 example.go@main:4]\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger.SetOutput(&buf)
			tt.doLog(logger)
			got := buf.String()
			if got != tt.want {
				t.Errorf("got:\n\t'%s'\nwanted:\n\t'%s'", got, tt.want)
			}
		})
	}
}

func TestWithDefaultLogLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := log.InitLog(log.WithDefaultLogLevel(log.Warn), log.WithOutput(&buf))
	logger.Print("foo")
	got := buf.String()
	want := "time = 2020-06-01T00:00:00Z00:00; level = warn; message = foo\n"
	if want != got {
		t.Errorf("got:\n\t'%s'\nwanted:\n\t'%s'", got, want)
	}
}

func TestLog_SetWatch(t *testing.T) {
	logger := log.InitLog(
		log.WithTimeFormat(time.RFC3339),
		log.WithTimer(func(format string) string { return testNow.Format(time.RFC3339) }),
		log.WithFormatter(&log.JsonFormatter{}),
	)
	var (
		a   string
		buf bytes.Buffer
	)
	logger.SetOutput(&buf)
	logger.SetWatch("a", &a)
	logger.Debug()
	a = "new"
	logger.Debug()
	a = "flex"
	b := "it"
	logger.SetWatch("b", &b)
	logger.Debug()
	logger.UnsetWatch("a")
	logger.Debug()
	got := buf.String()
	want := `{"time":"` + testNow.Format(time.RFC3339) + `","level":"debug","message":"","watch":{"a":""}}` + "\n" +
		`{"time":"` + testNow.Format(time.RFC3339) + `","level":"debug","message":"","watch":{"a":"new"}}` + "\n" +
		`{"time":"` + testNow.Format(time.RFC3339) + `","level":"debug","message":"","watch":{"a":"flex","b":"it"}}` + "\n" +
		`{"time":"` + testNow.Format(time.RFC3339) + `","level":"debug","message":"","watch":{"b":"it"}}` + "\n"
	if want != got {
		t.Errorf("got:\n\t'%s'\nwanted:\n\t'%s'", got, want)
	}

}

func TestLog_WithFields(t *testing.T) {
	logger := log.InitLog(
		log.WithTimer(func(format string) string { return testNow.Format(time.RFC3339) }),
		log.WithFormatter(&log.JsonFormatter{}),
		log.WithFieldOrder("time", "level", "message"),
	)
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	logger.Debug()
	logger.WithFields(log.InitFields().AddField("test", "case")).Debug()
	logger.Debug()
	got := buf.String()
	buf.Truncate(0)
	want := `{"time":"` + testNow.Format(time.RFC3339) + `","level":"debug","message":""}` + "\n" +
		`{"time":"` + testNow.Format(time.RFC3339) + `","level":"debug","message":"","test":"case"}` + "\n" +
		`{"time":"` + testNow.Format(time.RFC3339) + `","level":"debug","message":""}` + "\n"
	if want != got {
		t.Errorf("got:\n\t'%s'\nwanted:\n\t'%s'", got, want)
	}
}

func TestMultipleInstanceFields(t *testing.T) {
	var buf bytes.Buffer
	logger := log.InitLog(
		log.WithTimer(func(format string) string { return testNow.Format(time.RFC3339) }),
		log.WithFormatter(&log.JsonFormatter{}),
		log.WithFields(log.KV{"message", "always"}),
		log.WithFieldOrder("time", "level", "message"),
		log.WithOutput(&buf),
	)
	logger.Debug()
	logger.Debug("be testing")
	got := buf.String()
	want := `{"time":"` + testNow.Format(time.RFC3339) + `","level":"debug","message":["always",""]}` + "\n" +
		`{"time":"` + testNow.Format(time.RFC3339) + `","level":"debug","message":["always","be testing"]}` + "\n"
	if want != got {
		t.Errorf("got:\n\t'%s'\nwanted:\n\t'%s'", got, want)
	}
}

func TestLog_SetPrefix(t *testing.T) {
	var buf bytes.Buffer
	logger := log.InitLog(
		log.WithTimer(func(format string) string { return testNow.Format(time.RFC3339) }),
		log.WithFormatter(&log.JsonFormatter{}),
		log.WithFieldOrder("time", "level", "message"),
		log.WithOutput(&buf),
	)
	logger.SetPrefix("SET: ")
	want := `SET: {"time":"` + testNow.Format(time.RFC3339) + `","level":"debug","message":""}` + "\n" +
		`SET: {"time":"` + testNow.Format(time.RFC3339) + `","level":"debug","message":""}` + "\n"
	logger.Debug()
	logger.Debug()
	got := buf.String()
	if want != got {
		t.Errorf("got:\n\t'%s'\nwanted:\n\t'%s'", got, want)
	}
}

func TestLog_Prefix(t *testing.T) {
	var buf bytes.Buffer
	logger := log.InitLog(
		log.WithTimer(func(format string) string { return testNow.Format(time.RFC3339) }),
		log.WithFormatter(&log.JsonFormatter{}),
		log.WithFieldOrder("time", "level", "message"),
		log.WithOutput(&buf),
	)

	want := `{"time":"` + testNow.Format(time.RFC3339) + `","level":"debug","message":""}` + "\n" +
		`PREFIX: {"time":"` + testNow.Format(time.RFC3339) + `","level":"debug","message":""}` + "\n" +
		`{"time":"` + testNow.Format(time.RFC3339) + `","level":"debug","message":""}` + "\n"
	logger.Debug()
	logger.Prefix("PREFIX: ").Debug()
	logger.Debug()

	got := buf.String()
	if want != got {
		t.Errorf("got:\n\t'%s'\nwanted:\n\t'%s'", got, want)
	}
}

func TestLog_Panic(t *testing.T) {
	var buf bytes.Buffer
	logger := log.InitLog(
		log.WithFormatter(&log.JsonFormatter{}),
		log.WithFieldOrder("time", "level", "message"),
		log.WithOutput(&buf),
		log.WithPanic(func(v interface{}) { panic(v) }),
	)
	want := `{"time":"2020-06-01T00:00:00Z00:00","level":"panic","message":"","stack":["example.go@main:10","example.go@main:4"]}` + "\n"
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
		got := buf.String()
		if want != got {
			t.Errorf("got:\n\t'%s'\nwanted:\n\t'%s'", got, want)
		}
	}()
	logger.Panic()
}

func TestLog_Panicf(t *testing.T) {
	var buf bytes.Buffer
	logger := log.InitLog(
		log.WithFormatter(&log.JsonFormatter{}),
		log.WithFieldOrder("time", "level", "message"),
		log.WithOutput(&buf),
		log.WithPanic(func(v interface{}) { panic(v) }),
	)
	want := `{"time":"2020-06-01T00:00:00Z00:00","level":"panic","message":"2","stack":["example.go@main:10","example.go@main:4"]}` + "\n"
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
		got := buf.String()
		if want != got {
			t.Errorf("got:\n\t'%s'\nwanted:\n\t'%s'", got, want)
		}
	}()
	logger.Panicf("%d", 2)
}

func TestLog_Error_callStack(t *testing.T) {
	var buf bytes.Buffer
	logger := log.InitLog(log.WithCallStackMethod(log.GetCallStack), log.WithOutput(&buf))
	logger.Error(errors.New("foo"))
	got := buf.String()
	if !strings.Contains(got, "stack") {
		t.Error("expected a stack field, didn't get one in:\n" + got)
	}
}

// package log provides a highly customisable Log with:
// 	* sensible defaults
// 	* pre- and post-processing capabilities (interceptors)
// 	* spawnable sub-logs which can be modified/customised on their own
// 	* extensible formatting
// 	* custom timestamps
// 	* optional prefixing
// 	* custom exit behaviour on panic and fatal events.
// See the `Log` type for examples.
package log

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	maximumCallerDepth int = 25
	knownLoggerFrames  int = 5
)

var (
	defaultTimer       = func(format string) string { return time.Now().Format(format) }
	defaultStackMethod = GetCallStack
	defaultExit        = os.Exit
	defaultTimeFormat  = time.RFC3339
	defaultPanicMethod = func(v interface{}) { panic(v) }
)

// InitLog creates an instance of Log. If no options are presented, it is created with the following overrideable
// defaults (see examples for available options).
func InitLog(opts ...func(*Log)) *Log {
	l := &Log{
		levelNames:        make(map[LogLevel]string),
		defaultFieldNames: defaultFieldNames,
		defaultLogLevel:   Info,
		fields:            InitFields(),
		encoder:           &PlaintextFormatter{fieldSeparator: "; ", keyValSeparator: " = "},
		timer:             defaultTimer,
		minLevel:          Trace,
		minFrameLevel:     Error,
		fieldOrder:        []string{"time", "level", "message"},
		watches:           make(map[string]interface{}),
		timeFormat:        defaultTimeFormat,
		exit:              defaultExit,
		out:               os.Stdout,
		stackMethod:       defaultStackMethod,
		panicMethod:       defaultPanicMethod,
		mu:                sync.Mutex{},
	}
	for k, v := range defaultLevelNames {
		l.levelNames[k] = v
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// A Log holds configuration and provides functions to write log lines.
type Log struct {
	preProcessors     []func(*Fields)
	postProcessors    []func(*Fields)
	minLevel          LogLevel
	minFrameLevel     LogLevel
	enabledLevels     map[LogLevel]bool
	levelNames        map[LogLevel]string
	defaultFieldNames map[Label]string
	defaultLogLevel   LogLevel
	messageLabel      string
	out               io.Writer
	fields            *Fields
	encoder           Formatter
	timer             func(format string) string
	timeFormat        string
	fieldOrder        []string
	watches           map[string]interface{}
	exit              func(int)
	prefix            string
	stackMethod       func() []string
	panicMethod       func(v interface{})
	mu                sync.Mutex
}

// Log.Debug writes a log line of level `Debug` where the message field is a string made up of all values passed in.
func (l *Log) Debug(v ...interface{}) { l.Log(Debug, v...) }

// Log.Debugf formats a log message according to a format specifier and writes a log line of level `Debug`
func (l *Log) Debugf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Log(Debug, msg)
}

// Log.Error takes an error and some message parameters, adds an `error` field whose value is the error passed in and
// writes a log line of level `Error`
func (l *Log) Error(err error, v ...interface{}) {
	errorLog := l.WithFields(InitFields(KV{"error", err.Error()}))
	errorLog.Log(Error, v...)
}

// Log.Errorf formats a log message according to a format specifier, adds an `error` field whose value is the error passed in and
// writes a log line of level `Error`
func (l *Log) Errorf(err error, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	errorLog := l.WithFields(InitFields(KV{"error", err.Error()}))
	errorLog.Log(Error, msg)
}

// Log.Fatal takes an exit code and some message parameters, writes a log line of level 'Fatal' and exits the program.
func (l *Log) Fatal(code int, v ...interface{}) {
	l.Log(Fatal, v...)
	l.exit(code)
}

// Log.Logf formats a log message according to a format specifier, writes a log line of level 'Fatal' and exits the program.
func (l *Log) Fatalf(code int, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Log(Fatal, msg)
	l.exit(code)
}

// Log.Info writes a log line of level `Info` where the message field is a string made up of all values passed in.
func (l *Log) Info(v ...interface{}) { l.Log(Info, v...) }

// Log.Infof formats a log message according to a format specifier and writes a log line of level `Info`
func (l *Log) Infof(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Log(Info, l.messageLabel, msg)
}

// Log.Log  writes a log line with a level as provided.
func (l *Log) Log(level LogLevel, v ...interface{}) {
	if level < l.minLevel {
		return
	}
	l.preProcess(l.fields)
	logline := l.makeLogline(level, v...)
	l.writeLogline(logline)
	l.postProcess(logline)
}

func (l *Log) loggingError(err error, logline *Fields) {
	fmt.Fprintf(os.Stderr, "could not encode logline: %v\n", err)
	fmt.Fprintf(os.Stderr, "problem logline: %v\n", logline)
}

func (l *Log) makeLogline(level LogLevel, v ...interface{}) *Fields {
	logline := InitFields()

	logline.MergeFields(l.fields)
	logline.SetField(l.defaultFieldNames[TimestampLabel], l.timer(l.timeFormat))
	logline.SetField(l.defaultFieldNames[LevelLabel], l.levelNames[level])
	logline.AddField(l.defaultFieldNames[MessageLabel], fmt.Sprint(v...))
	if level >= l.minFrameLevel {
		logline.AddField(l.defaultFieldNames[StackLabel], l.stackMethod())
	}
	if len(l.watches) > 0 {
		logline.AddField(l.defaultFieldNames[WatchLabel], l.watchList())
	}
	logline.SetOrder(l.fieldOrder...)
	return logline
}

// Log.Panic writes a log line of level `Panic` where the message field is a string made up of all values passed in.
// It then panics, forwarding the parameters of this call.
func (l *Log) Panic(v ...interface{}) {
	l.Log(Panic, v...)
	l.panicMethod(fmt.Sprint(v...))
}

// Log.Panicf formats a log message according to a format specifier and writes a log line of level `Panic`
// It then panics, using the formatted message.
func (l *Log) Panicf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Log(Panic, msg)
	l.panicMethod(fmt.Sprint(msg))
}

func (l *Log) postProcess(logline *Fields) {
	for _, p := range l.postProcessors {
		p(logline)
	}
}

// Log.Prefix creates a new log instance with its own prefix

func (l *Log) Prefix(p string) *Log {
	prefixed := *l
	prefixed.SetPrefix(p)
	return &prefixed
}

func (l *Log) preProcess(logline *Fields) {
	for _, p := range l.preProcessors {
		p(logline)
	}
}

// Log.Print writes a log line with level set to the default, where the message field is a string made up of all values
// passed in.
func (l *Log) Print(v ...interface{}) { l.Log(l.defaultLogLevel, v...) }

// Log.Printf formats a log message according to a format specifier and writes a log line with level set the the default.
func (l *Log) Printf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Log(l.defaultLogLevel, l.messageLabel, msg)
}

// Log.SetOutput permanently changes the output writer for this instance.
func (l *Log) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

// Log.SetPrefix permanently changes the prefix for this instance.
func (l *Log) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}

// Log.SetWatch adds a pointer to a watch list with a label and will output the value of that reference whenever a log
// line is created. As many pointers as required can be added to the watch list. Any serializable type may be used.
func (l *Log) SetWatch(label string, pointer interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.watches[label] = pointer
}

// Log.Trace writes a log line of level `Trace` where the message field is a string made up of all values passed in.
func (l *Log) Trace(v ...interface{}) { l.Log(Trace, v...) }

// Log.Tracef formats a log message according to a format specifier and writes a log line of level `Trace`
func (l *Log) Tracef(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Log(Trace, l.messageLabel, msg)
}

// Log.UnsetWatch removes the specified element from the watch list.
func (l *Log) UnsetWatch(label string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.watches, label)
}

// Log.Warn writes a log line of level `Warn` where the message field is a string made up of all values passed in.
func (l *Log) Warn(v ...interface{}) { l.Log(Warn, v...) }

// Log.Warnf formats a log message according to a format specifier and writes a log line of level `Warn`
func (l *Log) Warnf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Log(Warn, l.messageLabel, msg)
}

func (l *Log) watchList() map[string]interface{} {
	l.mu.Lock()
	defer l.mu.Unlock()
	wl := make(map[string]interface{})
	vals, err := json.Marshal(&l.watches)
	if err != nil {
		l.Error(err, "unable to serialize watches")
	}
	err = json.Unmarshal(vals, &wl)
	if err != nil {
		l.Error(err, "unable to serialize watches")
	}
	return wl
}

// WithFields spawns a new log instance with additional *Fields as defined.
func (l *Log) WithFields(f *Fields) *Log {
	l.mu.Lock()
	defer l.mu.Unlock()
	fielded := *l
	fielded.fields = InitFields()
	fielded.fields.MergeFields(l.fields)
	fielded.fields.MergeFields(f)
	return &fielded
}

// Log.Writer returns the current output writer
func (l *Log) Writer() io.Writer { return l.out }

func (l *Log) writeLogline(logline *Fields) {
	bytes, err := l.encoder.Format(logline)
	line := append(bytes, []byte("\n")...)
	if err != nil {
		l.loggingError(err, logline)
		return
	}
	out := append([]byte(l.prefix), line...)
	if _, err := l.out.Write(out); err != nil {
		l.loggingError(err, logline)
		return
	}
}

// GetCallStack retrieves the call stack up to the point a given invocation was made.
func GetCallStack() []string {
	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(knownLoggerFrames, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	cs := make([]string, 0)
	for f, again := frames.Next(); again; f, again = frames.Next() {
		cs = append(cs, fmt.Sprintf("%s@%s:%d", f.Function, f.File, f.Line))
	}
	return cs
}

func DefaultTimer(timer func(string) string) bool {
	defaultTimer = timer
	return true
}

func DefaultStackMethod(stacker func() []string) bool {
	defaultStackMethod = stacker
	return true
}

func DefaultExit(exit func(int)) bool {
	defaultExit = exit
	return true
}

func DefaultPanic(pan func(v interface{})) bool {
	defaultPanicMethod = pan
	return true
}

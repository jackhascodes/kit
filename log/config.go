package log

import "io"

func WithCallStackMethod(stack func() []string) func(*Log) {
	return func(l *Log) {
		l.stackMethod = stack
	}
}

// WithDefaultLogLevel sets the level at which any "default" log actions log at (see `Log.Print`, `Log.Log`)
func WithDefaultLogLevel(level LogLevel) func(*Log) {
	return func(l *Log) {
		l.defaultLogLevel = level
	}
}

// WithExit allows the standard os.Exit called on `Log.Fatal` to be overridden or supplemented with prior action.
// The function passed as a parameter must have a single integer argument which should be an exit code.
func WithExit(exit func(int)) func(*Log) {
	return func(l *Log) {
		l.exit = exit
	}
}

// WithFieldOrder takes a variadic set of strings to set the order of output for fields.
// If fields are not represented in the parameter set, they will retain relative order but fields ordered explicitly
// with this option will take precedence.
// Duplicates are dealt with on a first-come-first-served basis, and any subsequent instance of a field name in the
// parameter set will be ignored.
func WithFieldOrder(o ...string) func(*Log) {
	return func(l *Log) {
		done := make(map[string]bool)
		new := make([]string, 0)
		for _, k := range o {
			done[k] = true
			new = append(new, k)
		}
		for _, k := range l.fieldOrder {
			if _, ok := done[k]; ok {
				continue
			}
			new = append(new, k)
		}
		l.fieldOrder = new
	}
}

// WithFields takes a variadic set of key-value pairs (`log.KV`) to use in every log line.
func WithFields(kv ...KV) func(*Log) {
	return func(l *Log) {
		l.fields.AddFields(kv...)
	}
}

// WithFormatter takes a Formatter interface (see log.Formatter) with which every log line will be formatted.
func WithFormatter(e Formatter) func(*Log) {
	return func(l *Log) {
		l.encoder = e
	}
}

// WithLevelNames sets the strings used to describe log levels. Any level not specifically defined will fall back to
// default (lowercase version of the constant name).
func WithLevelNames(levels map[LogLevel]string) func(*Log) {
	return func(l *Log) {
		for k, v := range levels {
			l.levelNames[k] = v
		}
	}
}

// WithMinLevel sets the minimum log level to be output. Any log call made with a level less than defined here is ignored.
// Ordering is defined by the `LogLevel` type constants.
func WithMinLevel(level LogLevel) func(*Log) {
	return func(l *Log) {
		l.minLevel = level
	}
}

// WithOutput defines where the log is written.
func WithOutput(o io.Writer) func(*Log) {
	return func(l *Log) {
		l.out = o
	}
}

func WithPanic(panicMethod func(interface{})) func(*Log) {
	return func(l *Log) {
		l.panicMethod = panicMethod
	}
}

// WithPrefix sets a string which will prefix all log lines.
func WithPrefix(p string) func(*Log) {
	return func(l *Log) {
		l.prefix = p
	}
}

// WithPreProcessors takes a variadic set of functions to be called against fields in a logline *BEFORE* any default
// fields are set (time, level, message).
func WithPreProcessors(p ...func(*Fields)) func(*Log) {
	return func(l *Log) {
		l.preProcessors = append(l.preProcessors, p...)
	}
}

// WithPostProcessors takes a variadic set of functions to be called against fields in a logline *AFTER* default fields
// are set (time, level, message).
func WithPostProcessors(p ...func(*Fields)) func(*Log) {
	return func(l *Log) {
		l.postProcessors = append(l.postProcessors, p...)
	}
}

// WithTimeFormat defines the format string to be used when writing the `time` field.
func WithTimeFormat(f string) func(*Log) {
	return func(l *Log) {
		l.timeFormat = f
	}
}

// WithTimer enables a custom time-keeper which takes a format string (defined by `WithTimeFormat()` or defaulted to
// `time.RFC3339`, performs a time operation and outputs a formatted string.
// Examples:
func WithTimer(t func(string) string) func(*Log) {
	return func(l *Log) {
		l.timer = t
	}
}

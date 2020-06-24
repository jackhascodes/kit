package log

type LogLevel int

const (
	Trace LogLevel = iota
	Debug
	Info
	Warn
	Error
	Fatal
	Panic
)

var defaultLevelNames = map[LogLevel]string{
	Error: "error",
	Fatal: "fatal",
	Info:  "info",
	Trace: "trace",
	Warn:  "warn",
	Debug: "debug",
	Panic: "panic",
}

package log

type Label int

const (
	LevelLabel Label = iota
	ErrorLabel
	MessageLabel
	TimestampLabel
	StackLabel
	WatchLabel
)

var defaultFieldNames = map[Label]string{
	LevelLabel:     "level",
	ErrorLabel:     "error",
	MessageLabel:   "message",
	TimestampLabel: "time",
	StackLabel:     "stack",
	WatchLabel:     "watch",
}

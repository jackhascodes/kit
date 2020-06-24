package log_test

import (
	"github.com/jackhascodes/kit/log"
)

func ExampleInitFields() {
	logger := log.InitLog()
	fields := log.InitFields(log.KV{"foo", "bar"}, log.KV{"baz", "qux"})
	logger.WithFields(fields).Debug("quux")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = debug; message = quux; foo = bar; baz = qux
}

func ExampleFields_AddField() {
	logger := log.InitLog()
	fields := log.InitFields(log.KV{"foo", "bar"})
	fields.AddField("baz", "qux")
	logger.WithFields(fields).Debug("quux")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = debug; message = quux; foo = bar; baz = qux
}

func ExampleFields_AddField_existing() {
	logger := log.InitLog()
	fields := log.InitFields(log.KV{"foo", "bar"})
	fields.AddField("baz", "qux")
	fields.AddField("baz", "gurp")
	logger.WithFields(fields).Debug("quux")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = debug; message = quux; foo = bar; baz = [qux gurp]
}

func ExampleFields_SetField() {
	logger := log.InitLog()
	fields := log.InitFields(log.KV{"foo", "bar"})
	fields.SetField("baz", "qux")
	logger.WithFields(fields).Debug("quux")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = debug; message = quux; foo = bar; baz = qux
}

func ExampleFields_SetField_existing() {
	logger := log.InitLog()
	fields := log.InitFields(log.KV{"foo", "bar"})
	fields.SetField("baz", "qux")
	fields.SetField("baz", "gurp")
	logger.WithFields(fields).Debug("quux")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = debug; message = quux; foo = bar; baz = gurp
}

func ExampleFields_SetField_override() {
	logger := log.InitLog()
	fields := log.InitFields(log.KV{"foo", "bar"})
	fields.AddField("baz", "qux")
	fields.SetField("baz", "gurp")
	logger.WithFields(fields).Debug("quux")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = debug; message = quux; foo = bar; baz = gurp
}

func ExampleFields_AddFields() {
	logger := log.InitLog()
	fields := log.InitFields(log.KV{"foo", "bar"})
	fields.AddFields(
		log.KV{"baz", "qux"},
		log.KV{"gurp", "flargle"},
	)
	logger.WithFields(fields).Debug("quux")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = debug; message = quux; foo = bar; baz = qux; gurp = flargle
}

func ExampleFields_AddFields_existing() {
	logger := log.InitLog()
	fields := log.InitFields(log.KV{"foo", "bar"})
	fields.AddFields(
		log.KV{"foo", "baz"},
		log.KV{"qux", "gurp"},
	)
	logger.WithFields(fields).Debug("quux")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = debug; message = quux; foo = [bar baz]; qux = gurp
}

func ExampleFields_AddFields_duplicated() {
	logger := log.InitLog()
	fields := log.InitFields(log.KV{"foo", "bar"})
	fields.AddFields(
		log.KV{"baz", "qux"},
		log.KV{"baz", "gurp"},
	)
	logger.WithFields(fields).Debug("quux")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = debug; message = quux; foo = bar; baz = [qux gurp]
}

func ExampleFields_SetOrder() {
	logger := log.InitLog()
	fields := log.InitFields(
		log.KV{"foo", "bar"},
		log.KV{"baz", "qux"},
	)
	fields.SetOrder("baz", "foo")
	logger.WithFields(fields).Debug("quux")
	// Output:
	// time = 2020-06-01T00:00:00Z00:00; level = debug; message = quux; baz = qux; foo = bar
}

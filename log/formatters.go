package log

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Formatter interface {
	Format(*Fields) ([]byte, error)
}

type JsonFormatter struct{}

func (j *JsonFormatter) Format(f *Fields) ([]byte, error) {
	var o = []byte("{")
	var parts = make([]string, 0)
	for _, k := range f.order {
		val, err := json.Marshal(f.values[k])
		if err != nil {
			return nil, err
		}
		parts = append(parts, "\""+k+"\":"+string(val))
	}
	o = append(o, []byte(strings.Join(parts, ","))...)
	o = append(o, []byte("}")...)
	return o, nil
}

type PlaintextFormatter struct {
	fieldSeparator  string
	keyValSeparator string
}

func (p *PlaintextFormatter) Format(f *Fields) ([]byte, error) {
	var o []string
	for _, k := range f.order {
		o = append(o, k+p.keyValSeparator+fmt.Sprintf("%v", f.values[k]))
	}
	return []byte(strings.Join(o, p.fieldSeparator)), nil
}

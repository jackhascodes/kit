package log

import (
	"errors"
	"fmt"
	"sync"
)

// TODO: Figure out how these will output in a deterministic manner
type Fields struct {
	values map[string]interface{}
	order  []string
	mu     sync.Mutex
}

type KV struct {
	Key   string
	Value interface{}
}

// InitFields initializes an instance of Fields, which can then be plugged into a
// Log instance.
func InitFields(keyValues ...KV) *Fields {
	values := make(map[string]interface{})
	order := make([]string, 0)
	for _, kv := range keyValues {
		values[kv.Key] = kv.Value
		order = append(order, kv.Key)
	}

	return &Fields{
		values: values,
		order:  order,
	}
}

// Fields.AddField allows a KV instance to be added to a Fields instance using simple key+value params.
// If a field key already exists, the existing value will be converted to a slice and the new value appended.
// If a field key does not already exist, it will be set as the last field in ordering.
func (f *Fields) AddField(key string, value interface{}) *Fields {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.values[key]; ok {
		multi := []interface{}{f.values[key]}
		multi = append(multi, value)
		value = multi
	}
	f.values[key] = value
	for _, v := range f.order {
		if v == key {
			return f
		}
	}
	f.order = append(f.order, key)
	return f
}

// Fields.Get gets the value of a field.
func (f *Fields) Get(key string) (interface{}, error) {
	if _, ok := f.values[key]; !ok {
		return nil, errors.New(fmt.Sprintf("field '%s' does not exist", key))
	}
	return f.values[key], nil
}

// Fields.SetField will override the existing value of a field key.
// If the field key does not currently exist, it will be added and set as the last field in ordering.
func (f *Fields) SetField(key string, value interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.values[key]; !ok {
		f.order = append(f.order, key)
	}
	f.values[key] = value
}

// Fields.AddFields allows the addition of multiple key value pairs (KV)
// If a field key already exists, the existing value will be converted to a slice and the new value appended.
// If a field key does not already exist, it will be set as the last field in ordering.
func (f *Fields) AddFields(keyValues ...KV) *Fields {
	for _, kv := range keyValues {
		f.AddField(kv.Key, kv.Value)
	}
	return f
}

// Fields.MergeFields adds the KVs from another Fields instance
func (f *Fields) MergeFields(fields *Fields) *Fields {
	for _, k := range fields.order {
		f.SetField(k, fields.values[k])
	}
	return f
}

// Fields.SetOrder re-orders keys based on their order in the params.
// Any keys not specified will retain their relative order at the end of the list.
func (f *Fields) SetOrder(o ...string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	done := make(map[string]bool)
	new := make([]string, 0)
	for _, k := range o {
		if _, ok := f.values[k]; !ok {
			continue
		}
		done[k] = true
		new = append(new, k)
	}
	for _, k := range f.order {
		if _, ok := done[k]; ok {
			continue
		}
		new = append(new, k)
	}
	f.order = new
}

// Fields.MarshalJSON overrides standard JSON marshalling, using the JsonFormatter
// to do the job instead. See `JsonFormatter`.
func (f *Fields) MarshalJSON() ([]byte, error) {
	format := &JsonFormatter{}
	return format.Format(f)
}

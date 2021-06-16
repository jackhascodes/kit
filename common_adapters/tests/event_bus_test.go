package tests

import (
	"errors"
	adapters "github.com/jackhascodes/kit/common_adapters"
	"testing"
)

type subscriber struct {
	handled int
}

func (s *subscriber) subscriber(label string) func(topic string, event interface{}) error {
	return func(topic string, event interface{}) error {
		s.handled++
		return nil
	}
}

func TestEventBus_Publish(t *testing.T) {
	bus := adapters.NewEventBus()
	sub := &subscriber{0}
	log := adapters.NewEventLogger()
	bus.Subscribe("*", log.Log)
	bus.Subscribe("a", sub.subscriber("a"))
	bus.Subscribe("a.b", sub.subscriber("a.b"))
	bus.Subscribe("a.*", sub.subscriber("a.*"))
	bus.Subscribe("*.b", sub.subscriber("*.b"))

	bus.Publish("a", 1)
	if sub.handled != 1 {
		t.Errorf("expected 1 got %d", sub.handled)
	}
	bus.Publish("a.b", 1)
	if sub.handled != 4 {
		t.Errorf("expected 4 got %d", sub.handled)
	}
	bus.Publish("b.b", 1)
	if sub.handled != 5 {
		t.Errorf("expected 5 got %d", sub.handled)
	}
	bus.Publish("c", errors.New("foo"))

}

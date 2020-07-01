package messagequeue

import (
	"testing"

	"github.com/jackhascodes/kit/messagequeue/msg"
)

func TestSubs(t *testing.T) {
	received_a := make([]*msg.Message, 0)
	received_b := make([]*msg.Message, 0)
	mq := InitMQ(Mock)

	mq.Connect()
	mq.Subscribe("test", func(message *msg.Message) {
		received_a = append(received_a, message)
	})
	mq.Subscribe("test", func(message *msg.Message) {
		received_b = append(received_b, message)
	})
	mq.Subscribe("apple", func(message *msg.Message) {
		received_b = append(received_b, message)
	})
	mq.Publish(&msg.Message{Topic: "test", Body: []byte("1")})
	mq.Publish(&msg.Message{Topic: "test", Body: []byte("2")})
	mq.Publish(&msg.Message{Topic: "test", Body: []byte("3")})
	mq.Publish(&msg.Message{Topic: "test", Body: []byte("4")})
	mq.Publish(&msg.Message{Topic: "test", Body: []byte("5")})
	mq.Unsubscribe("test")
	mq.Publish(&msg.Message{Topic: "test", Body: []byte("6")})
	mq.Publish(&msg.Message{Topic: "apple", Body: []byte("7")})
	mq.Close()

	want_a := 5
	want_b := 6
	got_a := len(received_a)
	got_b := len(received_b)
	if want_a != got_a {
		t.Errorf("expected %d messages received, got %d", want_a, got_a)
	}
	if want_b != got_b {
		t.Errorf("expected %d messages received, got %d", want_b, got_b)
	}
}

func TestQSubs(t *testing.T) {
	a_received := make([]*msg.Message, 0)
	b_received := make([]*msg.Message, 0)
	mq := InitMQ(Mock)
	mq.Connect()
	mq.QueueSubscribe("test", "a", func(message *msg.Message) {
		a_received = append(a_received, message)
	})
	mq.QueueSubscribe("test", "a", func(message *msg.Message) {
		b_received = append(b_received, message)
	})

	mq.QueueSubscribe("apple", "a", func(message *msg.Message) {
		b_received = append(b_received, message)
	})

	mq.Publish(&msg.Message{Topic: "test"})  // want a (1)
	mq.Publish(&msg.Message{Topic: "test"})  // want b (1)
	mq.Publish(&msg.Message{Topic: "test"})  // want a (2)
	mq.Publish(&msg.Message{Topic: "apple"}) // want b (2)
	mq.Unsubscribe("test")
	mq.Publish(&msg.Message{Topic: "test"})  // want none
	mq.Publish(&msg.Message{Topic: "apple"}) // want b (3)
	mq.Close()
	want_a := 2
	want_b := 3
	got_a := len(a_received)
	got_b := len(b_received)
	if want_a != got_a {
		t.Errorf("expected %d messages received on a, got %d", want_a, got_a)
	}
	if want_b != got_b {
		t.Errorf("expected %d messages received on b, got %d", want_b, got_b)
	}
}

package engines

import (
	"github.com/jackhascodes/kit/messagequeue/msg"
	natsio "github.com/nats-io/nats.go"
)

type Nats struct {
	conn     *natsio.Conn
	clientId string
	host     string
	subs     map[string]*natsio.Subscription
}

func InitNats(cfg *msg.Config) *Nats {
	return &Nats{
		clientId: cfg.ClientId,
		host:     cfg.Host,
	}
}

func (e *Nats) Connect() {
	e.conn, _ = natsio.Connect(e.host, natsio.Name(e.clientId))
}

func (e *Nats) Close() {
	for _, sub := range e.subs {
		sub.Unsubscribe()
	}
	e.conn.Close()
}

func (e *Nats) Publish(msg msg.Message) error {
	if msg.ReplyTopic == "" {
		return e.conn.Publish(msg.Topic, msg.Body)
	}
	return e.conn.PublishRequest(msg.Topic, msg.ReplyTopic, msg.Body)
}

func (e *Nats) Subscribe(topic string, h msg.Handler) error {
	sub, err := e.conn.Subscribe(topic, func(m *natsio.Msg) {
		h(msg.Message{
			Topic:      m.Subject,
			Body:       m.Data,
			ReplyTopic: m.Reply,
		})
	})
	if err != nil {
		return err
	}
	e.subs[topic] = sub
	return nil
}

func (e *Nats) QueueSubscribe(topic, queue string, h msg.Handler) error {
	sub, err := e.conn.QueueSubscribe(topic, queue, func(m *natsio.Msg) {
		h(msg.Message{
			Topic:      m.Subject,
			Body:       m.Data,
			ReplyTopic: m.Reply,
		})
	})

	if err != nil {
		return err
	}
	e.subs[topic] = sub
	return nil
}

func (e *Nats) Unsubscribe(topic string) {
	for _, sub := range e.subs {
			sub.Unsubscribe()
	}
}

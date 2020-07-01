package engines

import (
	"github.com/jackhascodes/kit/messagequeue/msg"
	rabbitio "github.com/streadway/amqp"
)

type Rabbit struct {
	conn     *rabbitio.Connection
	clientId string
	host     string
	subs     map[string]*rabbitio.Channel
}

func InitRabbit(cfg *msg.Config) *Rabbit {
	return &Rabbit{
		clientId: cfg.ClientId,
		host:     cfg.Host,
	}
}

func (e *Rabbit) Connect() {
	e.conn, _ = rabbitio.Dial(e.host)
}

func (e *Rabbit) Close() {
	for _, sub := range e.subs {
		sub.Close()
	}
	e.conn.Close()
}

func (e *Rabbit) Publish(msg *msg.Message) error {
	ch, _ := e.conn.Channel()
	defer ch.Close()
	e.declareExchange(ch, msg.Topic)

	if msg.ReplyTopic == "" {
		ch.Publish(msg.Topic, "", false, false, rabbitio.Publishing{Body: msg.Body})
	}
	return ch.Publish(msg.Topic, "", false, false, rabbitio.Publishing{Body: msg.Body, ReplyTo: msg.ReplyTopic})
}

func (e *Rabbit) Subscribe(topic string, h msg.Handler) error {
	return e.QueueSubscribe(topic, "", h)
}

func (e *Rabbit) QueueSubscribe(topic, queue string, h msg.Handler) error {
	ch, _ := e.conn.Channel()
	e.declareExchange(ch, topic)
	q, _ := e.declareQueue(ch, queue)
	e.bindQueue(ch, q, topic)
	e.subs[topic] = ch
	msgs, _ := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	go func() {
		for d := range msgs {
			h(&msg.Message{Topic: topic, Body: d.Body, ReplyTopic: d.ReplyTo})
		}
	}()
	return nil
}

func (e *Rabbit) Unsubscribe(topic string) {
	e.subs[topic].Close()
}

func (e *Rabbit) declareExchange(ch *rabbitio.Channel, name string) error {
	return ch.ExchangeDeclare(
		name,     // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
}

func (e *Rabbit) declareQueue(ch *rabbitio.Channel, name string) (rabbitio.Queue, error) {
	return ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

func (e *Rabbit) bindQueue(ch *rabbitio.Channel, q rabbitio.Queue, topic string) error {
	return ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		topic,  // exchange
		false,
		nil,
	)
}

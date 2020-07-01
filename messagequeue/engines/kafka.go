package engines

import (
	"context"
	"strings"
	"time"

	"github.com/jackhascodes/kit/messagequeue/msg"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
)

type Kafka struct {
	clientId string
	dialer   *kafka.Dialer
	hosts    []string
	logger   msg.Log
	pubs     map[string]*kafka.Writer
	subs     map[string]*kafka.Reader
}

func InitKafka(cfg *msg.Config) *Kafka {
	return &Kafka{
		clientId: cfg.ClientId,
		hosts:    strings.Split(cfg.Host, ","),
	}
}

func (e *Kafka) Connect() {
	// Connect to a server
	e.dialer = &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: e.clientId,
	}

}

func (e *Kafka) Close() {}

func (e *Kafka) Publish(msg *msg.Message) error {
	if _, ok := e.pubs[msg.Topic]; !ok {
		config := kafka.WriterConfig{
			Brokers:          e.hosts,
			Topic:            msg.Topic,
			Balancer:         &kafka.LeastBytes{},
			Dialer:           e.dialer,
			WriteTimeout:     10 * time.Second,
			ReadTimeout:      10 * time.Second,
			CompressionCodec: snappy.NewCompressionCodec(),
		}
		e.pubs[msg.Topic] = kafka.NewWriter(config)
	}
	return e.pubs[msg.Topic].WriteMessages(context.Background(), kafka.Message{
		Value: msg.Body,
	})
}

func (e *Kafka) Subscribe(topic string, h msg.Handler) error {
	config := kafka.ReaderConfig{
		Brokers:         e.hosts,
		GroupID:         e.dialer.ClientID,
		Topic:           topic,
		MinBytes:        10e3,            // 10KB
		MaxBytes:        10e6,            // 10MB
		MaxWait:         1 * time.Second, // Maximum amount of time to wait for new data to come when fetching batches of messages from kafka.
		ReadLagInterval: -1,
	}
	reader := kafka.NewReader(config)
	defer e.Unsubscribe(topic)
	e.subs[topic] = reader
	go func() {
		for {
			m, err := reader.ReadMessage(context.Background())
			if err != nil {
				continue
			}
			reply := getReplyTopic(m.Headers)
			h(&msg.Message{Topic: topic, Body: m.Value, ReplyTopic: reply})
		}
	}()
	return nil
}

func (e *Kafka) QueueSubscribe(topic, queue string, h msg.Handler) error {
	config := kafka.ReaderConfig{
		Brokers:         e.hosts,
		GroupID:         e.dialer.ClientID,
		Topic:           topic,
		MinBytes:        10e3,            // 10KB
		MaxBytes:        10e6,            // 10MB
		MaxWait:         1 * time.Second, // Maximum amount of time to wait for new data to come when fetching batches of messages from kafka.
		ReadLagInterval: -1,
	}
	reader := kafka.NewReader(config)
	defer e.Unsubscribe(topic)
	e.subs[topic] = reader
	go func() {
		for {
			m, err := reader.ReadMessage(context.Background())
			if err != nil {
				continue
			}
			reply := getReplyTopic(m.Headers)
			h(&msg.Message{Topic: topic, Body: m.Value, ReplyTopic: reply})
		}
	}()
	return nil
}

func (e *Kafka) Unsubscribe(topic string) {
	e.subs[topic].Close()
}

func getReplyTopic(headers []kafka.Header) string {
	for _, header := range headers {
		if header.Key == "ReplyTopic" {
			return string(header.Value)
		}
	}
	return ""
}

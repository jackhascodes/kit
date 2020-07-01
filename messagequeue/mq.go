package messagequeue

import (
	"log"

	"github.com/google/uuid"
	"github.com/jackhascodes/kit/messagequeue/engines"
	"github.com/jackhascodes/kit/messagequeue/msg"
)

type Engine interface {
	Connect()
	Close()
	Publish(msg *msg.Message) error
	Subscribe(string, msg.Handler) error
	QueueSubscribe(string, string, msg.Handler) error
	Unsubscribe(string)
}

// MQLog is a simple log wrapper to satisfy the Log interface used throughout as a default logger.
type MQLog struct{}

func (m *MQLog) Error(v ...interface{}) {
	log.Print(append([]interface{}{"MQ ERROR! "}, v...))
}

type MQ struct {
	engine Engine
}

type EngineType int

const (
	Mock EngineType = iota
	Nats
	Kafka
	Rabbit
)

func InitMQ(engine EngineType, opts ...func(*msg.Config)) *MQ {

	switch engine {
	case Mock:
		return InitMock(opts...)
	case Nats:
		return InitNats(opts...)
	case Kafka:
		return InitKafka(opts...)
	case Rabbit:
		return InitRabbit(opts...)
	default:
		return nil
	}
}

func InitKafka(opts ...func(*msg.Config)) *MQ {
	cfg := initConfig(opts...)
	mq := &MQ{engines.InitKafka(cfg)}
	mq.Connect()
	return mq
}

func InitMock(opts ...func(*msg.Config)) *MQ {
	cfg := initConfig(opts...)
	return &MQ{engines.InitMock(cfg)}
}

func InitNats(opts ...func(*msg.Config)) *MQ {
	cfg := initConfig(opts...)
	return &MQ{engines.InitNats(cfg)}
}

func InitRabbit(opts ...func(*msg.Config)) *MQ {
	cfg := initConfig(opts...)
	return &MQ{engines.InitRabbit(cfg)}
}

func initConfig(opts ...func(*msg.Config)) *msg.Config {
	cfg := &msg.Config{ClientId: uuid.New().String()}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

func WithClientId(clientId string) func(*msg.Config) {
	return func(cfg *msg.Config) {
		cfg.ClientId = clientId
	}
}

func WithClientIdPrefix(clientIdPrefix string) func(*msg.Config) {
	return func(cfg *msg.Config) {
		cfg.ClientId = clientIdPrefix + uuid.New().String()
	}
}

func WithHost(host string) func(*msg.Config) {
	return func(cfg *msg.Config) {
		cfg.Host = host
	}
}

func (m *MQ) Connect() {
	m.engine.Connect()
}

func (m *MQ) Publish(msg *msg.Message) error {
	return m.engine.Publish(msg)
}

func (m *MQ) Subscribe(topic string, h msg.Handler) {
	m.engine.Subscribe(topic, h)
}

func (m *MQ) QueueSubscribe(topic, queue string, h msg.Handler) {
	m.engine.QueueSubscribe(topic, queue, h)
}

func (m *MQ) Unsubscribe(topic string) {
	m.engine.Unsubscribe(topic)
}

func (m *MQ) Close() {
	m.engine.Close()
}

package engines

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/jackhascodes/kit/messagequeue/msg"
)

type Mock struct {
	ClientId     string
	Host         string
	Subs         map[string][]msg.Handler
	QSubs        map[string]map[string]*SimpleLB
	Conn         string
	PubQueue     map[string]*msg.Message
	LastReceived *msg.Message
	Worker       chan *msg.Message
	closed       chan bool
	mx           sync.RWMutex
	connstatus   chan bool
	closeRequest bool
	clearRequest bool
	clear        chan bool
}

func InitMock(cfg *msg.Config) *Mock {

	return &Mock{
		ClientId:   cfg.ClientId,
		Host:       cfg.Host,
		QSubs:      make(map[string]map[string]*SimpleLB),
		Subs:       make(map[string][]msg.Handler),
		PubQueue:   make(map[string]*msg.Message),
		Worker:     make(chan *msg.Message, 1),
		connstatus: make(chan bool, 1),
		closed:     make(chan bool, 1),
		clear:      make(chan bool, 1),
	}

}

func (e *Mock) Connect() {
	go func() {
		e.Conn = "connected"
		e.connstatus <- true
		for {
			message, ok := <-e.Worker
			if !ok {
				return
			}
			e.mx.Lock()
			if subs, ok := e.Subs[message.Topic]; ok {
				for _, sub := range subs {
					sub(message)
				}
			}
			if qsub, ok := e.QSubs[message.Topic]; ok {
				for _, lb := range qsub {
					if lb.Current == len(lb.Members) {
						lb.Current = 0
					}
					sub := lb.Members[lb.Current]
					sub(message)
					lb.Current++
				}
			}
			delete(e.PubQueue, message.Id)
			if e.closeRequest {
				if len(e.PubQueue) == 0 {
					e.close()
				}
			}

			if e.clearRequest {
				if len(e.PubQueue) == 0 {
					e.clear <- true
				}
			}
			e.mx.Unlock()
		}
	}()

}

func (e *Mock) Close() {
	e.closeRequest = true
	<-e.closed
}
func (e *Mock) close() {
	close(e.Worker)
	for topic, _ := range e.Subs {
		delete(e.Subs, topic)
	}
	for topic, _ := range e.QSubs {
		delete(e.QSubs, topic)
	}
	e.Conn = "sending closed"
	e.closed <- true

}

func (e *Mock) Publish(message *msg.Message) error {
	if e.Conn != "connected" {
		<-e.connstatus
	}
	if message.Topic == "error" {
		return errors.New("mock publish error")
	}
	message.Id = uuid.New().String()
	e.PubQueue[message.Id] = message
	e.Worker <- message
	return nil
}

func (e *Mock) Subscribe(topic string, h msg.Handler) error {
	e.mx.Lock()
	defer e.mx.Unlock()
	if _, ok := e.Subs[topic]; !ok {
		e.Subs[topic] = make([]msg.Handler, 0)
	}
	e.Subs[topic] = append(e.Subs[topic], h)
	return nil
}

func (e *Mock) QueueSubscribe(topic, queue string, h msg.Handler) error {
	if _, ok := e.QSubs[topic]; !ok {
		e.QSubs[topic] = make(map[string]*SimpleLB)
	}
	if _, ok := e.QSubs[topic][queue]; !ok {
		e.QSubs[topic][queue] = &SimpleLB{Members: make([]msg.Handler, 0), Current: 0}
	}
	e.QSubs[topic][queue].Members = append(e.QSubs[topic][queue].Members, h)
	return nil
}

func (e *Mock) Unsubscribe(topic string) {

	e.clearRequest = true
	<-e.clear
	e.mx.Lock()
	defer e.mx.Unlock()
	delete(e.Subs, topic)
	delete(e.QSubs, topic)
}

type SimpleLB struct {
	Members []msg.Handler
	Current int
}

package adapters

import (
	"github.com/sirupsen/logrus"
)
type EventLogger struct {
	log *logrus.Logger
}

func NewEventLogger() *EventLogger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	return &EventLogger{log}
}

func (l *EventLogger) Log(topic string, event interface{}) error {
	if err, ok := event.(error); ok {
		_ = l.Error(topic, err)
		return nil
	}
	_ = l.Info(topic, event)
	return nil
}

func (l *EventLogger) Info(topic string, event interface{}) error {
	l.log.WithFields(logrus.Fields{"detail": event}).Info(topic)
	return nil
}
func (l *EventLogger) Error(topic string, event interface{}) error {
	l.log.WithFields(logrus.Fields{"detail": event}).Error(topic)
	return nil
}

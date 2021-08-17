package common_adapters

import (
	"log"
	"regexp"
	"strings"
)

type EventBus struct {
	subscribers map[string]*SubscriberSet
}

type SubscriberSet struct {
	ephemeral []subscriptionHandler
	permanent []subscriptionHandler
}

type subscriptionHandler func(topic string, event interface{}) error

func NewEventBus() *EventBus {
	return &EventBus{subscribers: make(map[string]*SubscriberSet)}
}

func NewSubscriberSet() *SubscriberSet {
	return &SubscriberSet{make([]subscriptionHandler, 0), make([]subscriptionHandler, 0)}
}

func (e *EventBus) Subscribe(topic string, handlers ...subscriptionHandler) {
	if _, ok := e.subscribers[topic]; !ok {
		e.subscribers[topic] = NewSubscriberSet()
	}
	e.subscribers[topic].permanent = append(e.subscribers[topic].permanent, handlers...)
}

func (e *EventBus) SubscribeOnce(topic string, handlers ...subscriptionHandler) {
	if _, ok := e.subscribers[topic]; !ok {
		e.subscribers[topic] = NewSubscriberSet()
	}
	e.subscribers[topic].ephemeral = append(e.subscribers[topic].ephemeral, handlers...)
}

func (e *EventBus) Unsubscribe(topics ...string) {
	for _, topic := range topics {
		delete(e.subscribers, topic)
	}
}

func (e *EventBus) Publish(topic string, event interface{}) {
	topicMatches := e.findTopicMatches(topic)
	for _, key := range topicMatches {
		e.handleSubscribers(e.subscribers[key].permanent, topic, event, false)
		e.handleSubscribers(e.subscribers[key].ephemeral, topic, event, true)
		//e.subscribers[key].ephemeral = make([]subscriptionHandler, 0)
	}
}

func (e *EventBus) findTopicMatches(topic string) []string {
	matches := make([]string, 0)
	for key, _ := range e.subscribers {
		exp := prepareWilcardExpression(key)
		matcher := regexp.MustCompile(exp)
		if matcher.MatchString(topic) {
			matches = append(matches, key)
		}
	}
	return matches
}

func (e *EventBus) handleSubscribers(subs []subscriptionHandler, topic string, event interface{}, ephemeral bool) {
	for _, sub := range subs {
		err := sub(topic, event)
		if err != nil {
			log.Print("EVENT BUS ERROR", err.Error())
		}
	}
}

func prepareWilcardExpression(key string) string {
	key = strings.ReplaceAll(key, ".", "\\.")
	key = strings.ReplaceAll(key, "*", ".+")
	return "^" + key + "$"
}

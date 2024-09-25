// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bus

import (
	"context"
	"log"
)

type Handler func(event any)
type HandlerC func(c context.Context, event any)
type HandlerE func(event any) error
type HandlerCE func(c context.Context, event any) error

var subscribers map[string][]HandlerCE

func init() {
	subscribers = make(map[string][]HandlerCE)
}

func Subscribe(topic string, handler Handler) {
	SubscribeCE(topic, func(c context.Context, e any) error {
		handler(e)
		return nil
	})
}

func SubscribeE(topic string, handler HandlerE) {
	SubscribeCE(topic, func(c context.Context, e any) error {
		return handler(e)
	})
}

func SubscribeCE(topic string, handler HandlerCE) {
	_, exists := subscribers[topic]
	if !exists {
		subscribers[topic] = make([]HandlerCE, 0, 1)
	}
	subscribers[topic] = append(subscribers[topic], handler)
}

func Publish(topic string, event any) {
	if err := PublishE(topic, event); err != nil {
		log.Print("failed to publish event", topic, ":", err)
	}
}

func PublishC(c context.Context, topic string, event any) {
	if err := PublishCE(c, topic, event); err != nil {
		log.Print("failed to publish event ", topic, ": ", err)
	}
}

func PublishE(topic string, event any) error {
	return PublishCE(context.Background(), topic, event)
}

func PublishCE(c context.Context, topic string, event any) error {
	subs, ok := subscribers[topic]
	if !ok {
		return nil
	}
	for _, sub := range subs {
		if err := sub(c, event); err != nil {
			return err
		}
	}
	return nil
}

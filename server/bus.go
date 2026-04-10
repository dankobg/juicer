package server

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type bus struct {
	rdb         *redis.Client
	topics      []string
	subs        map[string]*redis.PubSub
	subMessages map[string]<-chan *redis.Message
}

func newBus(rdb *redis.Client) *bus {
	topics := []string{
		"ipc",
		"wsc.*",
	}

	bus := &bus{
		rdb:         rdb,
		topics:      topics,
		subs:        make(map[string]*redis.PubSub),
		subMessages: make(map[string]<-chan *redis.Message),
	}

	bus.subscribeToPubsub(context.Background())

	return bus
}

func (b *bus) subscribeToPubsub(ctx context.Context) {
	for _, topic := range b.topics {
		pubsub := b.rdb.PSubscribe(ctx, topic)
		b.subs[topic] = pubsub
		b.subMessages[topic] = pubsub.Channel()
	}
}

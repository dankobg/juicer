package ws

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
		"lobby*",
		"game.*",
		"gametv.*",
		"user.*",
		"conn.*",
	}

	bus := &bus{
		rdb:         rdb,
		topics:      topics,
		subs:        make(map[string]*redis.PubSub),
		subMessages: make(map[string]<-chan *redis.Message),
	}

	bus.subscribeToPubSub(context.Background())

	return bus
}

func (b *bus) subscribeToPubSub(ctx context.Context) {
	for _, topic := range b.topics {
		pubsub := b.rdb.PSubscribe(ctx, topic)
		b.subs[topic] = pubsub
		b.subMessages[topic] = pubsub.Channel()
	}
}

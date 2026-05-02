package server

import (
	"context"

	pb "github.com/dankobg/juicer/pb/proto/juicer"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
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

const useBinaryMessageFormat = false

func serializeMsg(msg *pb.Message) ([]byte, error) {
	if useBinaryMessageFormat {
		return proto.Marshal(msg)
	}

	return protojson.Marshal(msg)
}

func deserializeMsg(bb []byte, msg *pb.Message) error {
	if useBinaryMessageFormat {
		return proto.Unmarshal(bb, msg)
	}

	return protojson.Unmarshal(bb, msg)
}

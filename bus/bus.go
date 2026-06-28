package bus

import (
	"context"

	pb "github.com/dankobg/juicer/pb/proto/juicer"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Bus struct {
	rdb         *redis.Client
	topics      []string
	Subs        map[string]*redis.PubSub
	SubMessages map[string]<-chan *redis.Message
}

func NewBus(rdb *redis.Client) *Bus {
	topics := []string{
		"ipc",
		"wsc.*",
	}

	bus := &Bus{
		rdb:         rdb,
		topics:      topics,
		Subs:        make(map[string]*redis.PubSub),
		SubMessages: make(map[string]<-chan *redis.Message),
	}

	bus.subscribeToPubsub(context.Background())

	return bus
}

func (b *Bus) Publish(ctx context.Context, channel string, message any) error {
	return b.rdb.Publish(ctx, channel, message).Err()
}

func (b *Bus) subscribeToPubsub(ctx context.Context) {
	for _, topic := range b.topics {
		pubsub := b.rdb.PSubscribe(ctx, topic)
		b.Subs[topic] = pubsub
		b.SubMessages[topic] = pubsub.Channel(redis.WithChannelSize(10_000))
		// i have issues here, increasing ^ buffer won't solve the problem, just hides it
		// redis: pubsub.go:793: channel is full for 1m0s (message is dropped)
	}
}

func (b *Bus) Close() {
	for _, sub := range b.Subs {
		if sub != nil {
			_ = sub.Close()
		}
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

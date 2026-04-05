package server

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/dankobg/juicer/pb/proto/juicer"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/encoding/protojson"
)

func (a *ApiHandler) Test() error {
	return nil
}

func (a *ApiHandler) subscribeToPubsub(ctx context.Context) {
	topics := []string{
		"ipc",
		"wsc.*",
	}

	for _, topic := range topics {
		pubsub := a.Rdb.PSubscribe(ctx, topic)
		a.bus.subs[topic] = pubsub
		a.bus.subMessages[topic] = pubsub.Channel()
	}
}

func (a *ApiHandler) PubsubProcess(ctx context.Context) {
	a.Log.Info("gameserver pubsub started")

	for {
		select {
		case msg := <-a.bus.subMessages["ipc"]:
			a.handlePubsubRecvIPCMessage(msg)
		case msg := <-a.bus.subMessages["wsc."]:
			a.handlePubsubRecvLobbyMessage(msg)

		case <-ctx.Done():
			a.Log.Debug("gameserver pubsub ctx done")
			return
		}
	}
}

func (a *ApiHandler) handlePubsubRecvIPCMessage(msg *redis.Message) {
	m := &pb.Message{}
	if err := protojson.Unmarshal([]byte(msg.Payload), m); err != nil {
		a.Log.Error("protojson.Unmarshal IPC Message")
		return
	}

	switch m.GetEvent().(type) {
	case *pb.Message_RequestInitialChannels:
		requestInitialChannels := m.GetRequestInitialChannels()

		initialChannelsReplyMsg := &pb.Message{
			Event: &pb.Message_InitialChannels{InitialChannels: &pb.InitialChannelsReply{
				InitialChannels: []string{"loby", "game.123", "gametv.456"},
			}},
		}

		initialChannelsReplyMsgBytes, err := protojson.Marshal(initialChannelsReplyMsg)
		if err != nil {
			a.Log.Error("protojson marshal Message_InitialChannels", slog.String("client_id", requestInitialChannels.ClientId), slog.Any("error", err))
			return
		}
		topic := "reply-initial-channels." + requestInitialChannels.ClientId + "." + requestInitialChannels.ConnId
		if err := a.Rdb.Publish(context.Background(), topic, initialChannelsReplyMsgBytes).Err(); err != nil {
			a.Log.Error("hub publish Message_InitialChannels", slog.String("client_id", requestInitialChannels.ClientId), slog.String("topic", "ipc"), slog.Any("error", err))
			return
		}
	}
}

func (a *ApiHandler) handlePubsubRecvLobbyMessage(msg *redis.Message) {
	fmt.Println("gameserver handlePubsubRecvLobbyMessage", msg)
}

package server

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	pb "github.com/dankobg/juicer/pb/proto/juicer"
	"github.com/dankobg/juicer/ws"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/encoding/protojson"
)

type clientAuthInfo struct {
	clientID  string
	connID    string
	authState ws.ClientAuthState
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
			a.onIPCMsg(msg)
		case msg := <-a.bus.subMessages["wsc.*"]:
			a.onWSCMsg(msg)

		case <-ctx.Done():
			a.Log.Debug("gameserver pubsub ctx done")
			return
		}
	}
}

func (a *ApiHandler) onIPCMsg(m *redis.Message) {
	msg := &pb.Message{}
	if err := protojson.Unmarshal([]byte(m.Payload), msg); err != nil {
		a.Log.Error("protojson.Unmarshal IPC Message")
		return
	}

	switch msg.GetEvent().(type) {
	case *pb.Message_RequestInitialChannels:
		data := msg.GetRequestInitialChannels()

		initialChannelsReplyMsg := &pb.Message{
			Event: &pb.Message_InitialChannels{InitialChannels: &pb.InitialChannels{
				Channels: []string{"lobby", "lobby.chat"},
			}},
		}

		initialChannelsReplyMsgBytes, err := protojson.Marshal(initialChannelsReplyMsg)
		if err != nil {
			a.Log.Error("protojson marshal Message_InitialChannels", slog.String("client_id", data.ClientId), slog.Any("error", err))
			return
		}
		topic := "reply-initial-channels." + data.ClientId + "." + data.ConnId
		if err := a.Rdb.Publish(context.Background(), topic, initialChannelsReplyMsgBytes).Err(); err != nil {
			a.Log.Error("hub publish Message_InitialChannels", slog.String("client_id", data.ClientId), slog.String("topic", "ipc"), slog.Any("error", err))
			return
		}

	case *pb.Message_RequestChannelsInfo:
		data := msg.GetRequestChannelsInfo()
		_ = data
	}
}

func (a *ApiHandler) onWSCMsg(m *redis.Message) {
	msg := &pb.Message{}
	if err := protojson.Unmarshal([]byte(m.Payload), msg); err != nil {
		a.Log.Error("protojson.Unmarshal WSC Message")
		return
	}

	clientAuthInfo, err := extractWSCTopicParts(m.Channel)
	if err != nil {
		a.Log.Error("extractWSCTopicParts", slog.String("channel", m.Channel), slog.String("pattern", m.Pattern), slog.String("payload", m.Payload), slog.Any("error", err))
		return
	}

	switch msg.GetEvent().(type) {
	case *pb.Message_Test:
		a.handleWSCTestMsg(clientAuthInfo, msg.GetTest())
	}
}

func (a *ApiHandler) handleWSCTestMsg(authInfo clientAuthInfo, data *pb.Test) {
	xxx := &pb.Message{Event: &pb.Message_Test{Test: &pb.Test{Message: strings.ToUpper(data.Message)}}}
	b, _ := protojson.Marshal(xxx)

	topic := "user." + authInfo.clientID
	// topic := "conn." + authInfo.connID
	// topic := "lobby.chat"
	a.Rdb.Publish(context.Background(), topic, b)
}

// extractWSCTopicParts extracts the client_id and auth_state
func extractWSCTopicParts(topic string) (clientAuthInfo, error) {
	parts := strings.Split(topic, ".")
	if len(parts) != 4 {
		return clientAuthInfo{}, fmt.Errorf("invalid parts length, expected 4, got: %d", len(parts))
	}
	clientID, connID, authStateStr := parts[1], parts[2], parts[3]
	if !(authStateStr == "0" || authStateStr == "1") {
		return clientAuthInfo{}, fmt.Errorf("invalid parts auth_state. must be 0 or 1")
	}
	authState := ws.ClientGuest
	if authStateStr == "1" {
		authState = ws.ClientAuth
	}
	return clientAuthInfo{
		clientID:  clientID,
		connID:    connID,
		authState: authState,
	}, nil
}

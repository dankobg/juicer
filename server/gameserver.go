package server

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	pb "github.com/dankobg/juicer/pb/proto/juicer"
	"github.com/dankobg/juicer/ws"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/encoding/protojson"
)

type clientAuthInfo struct {
	userID    string
	connID    string
	authState ws.ClientAuthState
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
	case *pb.Message_Heartbeat:
		a.handleIPCHeartbeatMsg(msg.GetHeartbeat())

	case *pb.Message_LeaveTab:
		a.handleIPCLeaveTabMsg(msg.GetLeaveTab())

	case *pb.Message_LeaveSite:
		a.handleIPCLeaveSiteMsg(msg.GetLeaveSite())

	case *pb.Message_RequestInitialChannels:
		a.handleIPCRequestInitialChannelsMsg(msg.GetRequestInitialChannels())

	case *pb.Message_RequestChannelsInfo:
		a.handleIPCRequestChannelsInfoMsg(msg.GetRequestChannelsInfo())
	}
}

func (a *ApiHandler) handleIPCHeartbeatMsg(data *pb.Heartbeat) {

}

func (a *ApiHandler) handleIPCLeaveTabMsg(data *pb.LeaveTab) {

}

func (a *ApiHandler) handleIPCLeaveSiteMsg(data *pb.LeaveSite) {
}

func (a *ApiHandler) handleIPCRequestInitialChannelsMsg(data *pb.RequestInitialChannels) {
	initialChannelsReplyMsg := &pb.Message{
		Event: &pb.Message_InitialChannels{InitialChannels: &pb.InitialChannels{
			Channels: []string{"lobby", "lobby.chat"},
		}},
	}

	initialChannelsReplyMsgBytes, err := protojson.Marshal(initialChannelsReplyMsg)
	if err != nil {
		a.Log.Error("protojson marshal Message_InitialChannels", slog.String("user_id", data.UserId), slog.Any("error", err))
		return
	}

	topic := "reply-initial-channels." + data.UserId + "." + data.ConnId
	if err := a.bus.rdb.Publish(context.Background(), topic, initialChannelsReplyMsgBytes).Err(); err != nil {
		a.Log.Error("hub publish Message_InitialChannels", slog.String("user_id", data.UserId), slog.String("topic", "ipc"), slog.Any("error", err))
		return
	}
}

func (a *ApiHandler) handleIPCRequestChannelsInfoMsg(data *pb.RequestChannelsInfo) {

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
	case *pb.Message_Echo:
		a.handleWSCEchoMsg(clientAuthInfo, msg.GetEcho())
	case *pb.Message_SeekGame:
		a.handleWSCSeekGameMsg(clientAuthInfo, msg.GetSeekGame())
	case *pb.Message_CancelSeekGame:
		a.handleWSCCancelSeekGameMsg(clientAuthInfo, msg.GetCancelSeekGame())
	}
}

func (a *ApiHandler) handleWSCEchoMsg(authInfo clientAuthInfo, data *pb.Echo) {
	xxx := &pb.Message{Event: &pb.Message_Echo{Echo: &pb.Echo{Message: strings.ToUpper(data.Message)}}}
	b, _ := protojson.Marshal(xxx)

	topic := "user." + authInfo.userID
	// topic := "conn." + authInfo.connID
	// topic := "lobby.chat"
	a.bus.rdb.Publish(context.Background(), topic, b)
}

func (a *ApiHandler) handleWSCSeekGameMsg(authInfo clientAuthInfo, data *pb.SeekGame) {
	if data == nil {
		return
	}

	ctx := context.Background()

	if _, err := a.bus.rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		key := fmt.Sprintf("seek_game.%d.%d-%d", authInfo.authState, data.GetTimeControl().Clock.Seconds, data.GetTimeControl().Increment.Seconds)
		if err := p.ZAdd(ctx, key, redis.Z{Member: authInfo.userID, Score: float64(time.Now().UnixNano())}).Err(); err != nil {
			a.Log.Error("SeekGame add to queue", slog.String("user_id", authInfo.userID), slog.String("auth_state", authInfo.authState.String()), slog.Any("error", err))
		}

		// 	if err := p.HSet(ctx, "clients_seeking", clientID, key).Err(); err != nil {
		// 		h.log.Error("seek_game add seek key for client", slog.String("user_id", clientID), slog.String("auth_state", authState.String()), slog.Any("error", err))
		// 		return err
		// 	}

		if err := p.Publish(ctx, key, authInfo.userID).Err(); err != nil {
			a.Log.Error("SeekGame publish joined", slog.String("user_id", authInfo.userID), slog.String("auth_state", authInfo.authState.String()), slog.Any("error", err))
			return err
		}

		return nil
	}); err != nil {
		a.Log.Error("SeekGame pipeline", slog.String("user_id", authInfo.userID), slog.String("auth_state", authInfo.authState.String()), slog.Any("error", err))
	}

	// h.broadcastHubInfoToClient(uuid.MustParse(clientID))
}

func (a *ApiHandler) handleWSCCancelSeekGameMsg(authInfo clientAuthInfo, data *pb.CancelSeekGame) {
	fmt.Println(data, "CANCEL_SEEK")
}

// extractWSCTopicParts extracts the user_id, conn_id and auth_state
func extractWSCTopicParts(topic string) (clientAuthInfo, error) {
	parts := strings.Split(topic, ".")
	if len(parts) != 4 {
		return clientAuthInfo{}, fmt.Errorf("invalid parts length, expected 4, got: %d", len(parts))
	}

	clientID, connID, authStateStr := parts[1], parts[2], parts[3]
	if authStateStr != "0" && authStateStr != "1" {
		return clientAuthInfo{}, fmt.Errorf("invalid parts auth_state. must be 0 or 1")
	}

	authState := ws.ClientGuest
	if authStateStr == "1" {
		authState = ws.ClientAuth
	}

	return clientAuthInfo{
		userID:    clientID,
		connID:    connID,
		authState: authState,
	}, nil
}

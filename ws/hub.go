package ws

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	pb "github.com/dankobg/juicer/pb/proto/juicer"
	"github.com/dankobg/juicer/persistence"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/encoding/protojson"
)

type Hub struct {
	ClientConnected    chan *client
	ClientDisconnected chan *client
	clientsByID        map[uuid.UUID]map[*client]struct{}
	clientsByConnID    map[uuid.UUID]*client
	clientChannels     map[*client][]Channel
	channels           map[Channel]map[*client]struct{}
	mu                 *sync.Mutex
	subs               map[string]*redis.PubSub
	subMessages        map[string]<-chan *redis.Message
	broadcastConn      chan ConnMessage
	broadcastUser      chan UserMessage
	broadcastChannel   chan ChannelMessage
	rdb                *redis.Client
	log                *slog.Logger
}

func NewHub(persistor persistence.Persistor, rdb *redis.Client, logger *slog.Logger) *Hub {
	hub := &Hub{
		ClientConnected:    make(chan *client),
		ClientDisconnected: make(chan *client),
		clientsByID:        make(map[uuid.UUID]map[*client]struct{}),
		clientsByConnID:    make(map[uuid.UUID]*client),
		clientChannels:     make(map[*client][]Channel),
		channels:           make(map[Channel]map[*client]struct{}),
		mu:                 &sync.Mutex{},
		subs:               make(map[string]*redis.PubSub),
		subMessages:        make(map[string]<-chan *redis.Message),
		broadcastConn:      make(chan ConnMessage, 100),
		broadcastUser:      make(chan UserMessage, 100),
		broadcastChannel:   make(chan ChannelMessage, 100),
		rdb:                rdb,
		log:                logger,
	}

	hub.subscribeToPubsub(context.Background())

	return hub
}

func (h *Hub) subscribeToPubsub(ctx context.Context) {
	topics := []string{
		"lobby*",
		"game.*",
		"gametv.*",
		"user.*",
		"conn.*",
	}

	for _, topic := range topics {
		pubsub := h.rdb.PSubscribe(ctx, topic)
		h.subs[topic] = pubsub
		h.subMessages[topic] = pubsub.Channel()
	}
}

func (h *Hub) listenForRedisRequests() {
	for {
		//
	}
}

// Run starts the pubsub and machmaking, as well as broadcast events
func (h *Hub) Run(ctx context.Context) error {
	defer func() {
		if r := recover(); r != nil {
			h.log.Info("hub recovered", slog.Any("recover", r))
		}
	}()

	h.log.Info("hub is running")

	go h.PubsubProcess(ctx)

loop:
	for {
		select {
		case c := <-h.ClientConnected:
			h.onClientConnected(c)
		case c := <-h.ClientDisconnected:
			h.onClientDisconnected(c)
		case m, ok := <-h.broadcastConn:
			if !ok {
				continue
			}
			h.onBroadcastConn(m)
		case m, ok := <-h.broadcastUser:
			if !ok {
				continue
			}
			h.onBroadcastUser(m)
		case m, ok := <-h.broadcastChannel:
			if !ok {
				continue
			}
			h.onBroadcastChannel(m)
		case <-ctx.Done():
			break loop
		}
	}

	return nil
}

func (h *Hub) Stop() {
	for _, sub := range h.subs {
		_ = sub.Close()
	}
}

// processClientWebsocketMessage publishes client websocket message to pubsub
func (h *Hub) processClientWebsocketMessage(client *client, msg []byte) error {
	topic := fmt.Sprintf("wsc.%s.%s.%d", client.id, client.connID, client.authState)

	if err := h.rdb.Publish(context.Background(), topic, msg).Err(); err != nil {
		h.log.Error("hub publish msg from websocket", slog.String("client_id", client.id.String()), slog.String("conn_id", client.connID.String()), slog.String("topic", topic), slog.Any("error", err))
	}

	return nil
}

func (h *Hub) onClientConnected(client *client) {
	h.log.Debug("client connected", slog.String("client_id", client.id.String()), slog.String("auth_state", client.authState.String()))

	h.addClient(client)
	h.requestChannelsInfo(client)

	clientConnectedMsg := &pb.Message{
		Event: &pb.Message_ClientConnected{ClientConnected: &pb.ClientConnected{Id: client.id.String()}},
	}

	clientConnectedMsgBytes, err := protojson.Marshal(clientConnectedMsg)
	if err != nil {
		h.log.Error("protojson marshal Message_ClientConnected", slog.String("client_id", client.id.String()), slog.Any("error", err))
	} else {
		if err := h.rdb.Publish(context.Background(), "ipc", clientConnectedMsgBytes).Err(); err != nil {
			h.log.Error("hub publish Message_ClientConnected", slog.String("client_id", client.id.String()), slog.String("topic", "ipc"), slog.Any("error", err))
		}
	}
}

func (h *Hub) onClientDisconnected(client *client) {
	h.log.Debug("client disconnected", slog.String("client_id", client.id.String()), slog.String("auth_state", client.authState.String()))
	h.removeClient(client)

	clientDisconnectedMsg := &pb.Message{
		Event: &pb.Message_ClientDisconnected{ClientDisconnected: &pb.ClientDisconnected{Id: client.id.String()}},
	}

	clientDisconnectedMsgBytes, err := protojson.Marshal(clientDisconnectedMsg)
	if err != nil {
		h.log.Error("protojson marshal Message_ClientDisconnected", slog.String("client_id", client.id.String()), slog.Any("error", err))
	} else {
		if err := h.rdb.Publish(context.Background(), "ipc", clientDisconnectedMsgBytes).Err(); err != nil {
			h.log.Error("hub publish Message_ClientDisconnected", slog.String("client_id", client.id.String()), slog.String("topic", "ipc"), slog.Any("error", err))
		}
	}
}

func (h *Hub) onBroadcastConn(connMsg ConnMessage) {
	h.log.Info("broadcasting to conn", slog.String("conn_id", connMsg.connID.String()), slog.String("msg", string(connMsg.msg)))

	c, ok := h.clientsByConnID[connMsg.connID]
	if !ok {
		h.log.Debug("broadcasting to conn, conn not found", slog.String("conn_id", connMsg.connID.String()))
		return
	}
	select {
	case c.outMsg <- connMsg.msg:
	default:
		h.removeClient(c)
	}
}

func (h *Hub) onBroadcastUser(clientMsg UserMessage) {
	h.log.Info("broadcasting to client", slog.String("client_id", clientMsg.userID.String()), slog.String("msg", string(clientMsg.msg)))

	for c := range h.clientsByID[clientMsg.userID] {
		canSend := true

		if clientMsg.channel != nil {
			canSend = false

			for _, channel := range c.channels {
				if strings.HasPrefix(clientMsg.channel.String(), channel.String()) {
					canSend = true
					break
				}
			}
		}

		if !canSend {
			continue
		}

		select {
		case c.outMsg <- clientMsg.msg:
		default:
			h.removeClient(c)
		}
	}
}

func (h *Hub) onBroadcastChannel(channelMsg ChannelMessage) {
	h.log.Info("broadcasting to channel", slog.String("channel", channelMsg.channel.String()), slog.String("msg", string(channelMsg.msg)))

	for c := range h.channels[channelMsg.channel] {
		select {
		case c.outMsg <- channelMsg.msg:
		default:
			h.removeClient(c)
		}
	}
}

func (h *Hub) addClient(c *client) {
	h.mu.Lock()

	if h.clientsByID[c.id] == nil {
		h.clientsByID[c.id] = make(map[*client]struct{})
	}
	h.clientsByID[c.id][c] = struct{}{}
	h.clientsByConnID[c.connID] = c
	h.clientChannels = make(map[*client][]Channel)

	h.clientChannels[c] = make([]Channel, 0)
	for _, clientChannel := range c.channels {
		if h.channels[clientChannel] == nil {
			h.channels[clientChannel] = make(map[*client]struct{})
		}

		h.channels[clientChannel][c] = struct{}{}
		h.clientChannels[c] = append(h.clientChannels[c], clientChannel)
	}

	h.mu.Unlock()

	h.log.Info("client added", slog.String("client_id", c.id.String()), slog.String("conn_id", c.connID.String()), slog.String("auth_state", c.authState.String()))
}

func (h *Hub) removeClient(c *client) {
	close(c.outMsg)

	h.mu.Lock()

	for _, clientChannel := range h.clientChannels[c] {
		delete(h.channels[clientChannel], c)

		if len(h.channels[clientChannel]) == 0 {
			delete(h.channels, clientChannel)
		}
	}

	delete(h.clientChannels, c)
	delete(h.clientsByConnID, c.connID)

	if len(h.clientsByID[c.id]) == 1 {
		delete(h.clientsByID, c.id)

		// send backend that user has left the site, so it can do things like:
		// cancel seeks, inform game members that their opponent left etc.
	} else {
		delete(h.clientsByID[c.id], c)
	}

	h.mu.Unlock()

	h.log.Info("client removed", slog.String("client_id", c.id.String()), slog.String("conn_id", c.connID.String()), slog.String("auth_state", c.authState.String()))

}

func (h *Hub) RequestInitialChannels(ctx context.Context, client *client) ([]string, error) {
	h.log.Info("requesting initial channels", slog.String("client_id", client.id.String()), slog.String("conn_id", client.connID.String()), slog.String("auth_state", client.authState.String()))

	topic := "reply-initial-channels." + client.id.String() + "." + client.connID.String()
	sub := h.rdb.Subscribe(ctx, topic)
	defer func() {
		_ = sub.Close()
	}()

	requestInitialChannelsMsg := &pb.Message{
		Event: &pb.Message_RequestInitialChannels{RequestInitialChannels: &pb.RequestInitialChannels{
			ClientId: client.id.String(),
			ConnId:   client.connID.String(),
			Path:     client.query.Get("path"),
		}},
	}

	requestInitialChannelsMsgBytes, err := protojson.Marshal(requestInitialChannelsMsg)
	if err != nil {
		h.log.Error("protojson marshal Message_InitialChannels", slog.String("client_id", client.id.String()), slog.Any("error", err))
	} else {
		if err := h.rdb.Publish(context.Background(), "ipc", requestInitialChannelsMsgBytes).Err(); err != nil {
			h.log.Error("hub publish Message_RequestInitialChannels", slog.String("client_id", client.id.String()), slog.String("topic", "ipc"), slog.Any("error", err))
		}
	}

	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		h.log.Error("hub recv reply initial-channels", slog.String("client_id", client.id.String()), slog.String("topic", "ipc"), slog.Any("error", err))
		return nil, fmt.Errorf("failed to receive initial channels reply: %w", err)
	}

	m := &pb.Message{}
	if err := protojson.Unmarshal([]byte(msg.Payload), m); err != nil {
		h.log.Error("protojson.Unmarshal reply initial-channels message")
		return nil, fmt.Errorf("protojson.Unmarshal Message_InitialChannels: %w", err)
	}

	initialChannels := m.GetInitialChannels().GetChannels()

	return initialChannels, nil
}

func (h *Hub) requestChannelsInfo(client *client) {
	channels := make([]string, len(client.channels))
	for i, channel := range client.channels {
		channels[i] = channel.String()
	}

	requestChannelsInfoMsg := &pb.Message{
		Event: &pb.Message_RequestChannelsInfo{
			RequestChannelsInfo: &pb.RequestChannelsInfo{
				ClientId: client.id.String(),
				ConnId:   client.connID.String(),
				Channels: channels,
			},
		},
	}

	requestChannelsInfoMsgBytes, err := protojson.Marshal(requestChannelsInfoMsg)
	if err != nil {
		h.log.Error("protojson marshal Message_RequestChannelsInfo", slog.String("client_id", client.id.String()), slog.Any("error", err))
	} else {
		if err := h.rdb.Publish(context.Background(), "ipc", requestChannelsInfoMsgBytes).Err(); err != nil {
			h.log.Error("hub publish Message_RequestChannelsInfo", slog.String("client_id", client.id.String()), slog.String("topic", "ipc"), slog.Any("error", err))
		}
	}
}

func (h *Hub) PubsubProcess(ctx context.Context) {
	h.log.Info("hub pubsub started")

	for {
		select {
		case msg := <-h.subMessages["lobby*"]:
			h.handlePubsubRecvLobbyMessage(msg)
		case msg := <-h.subMessages["game.*"]:
			h.handlePubsubRecvGameMessage(msg)
		case msg := <-h.subMessages["gametv.*"]:
			h.handlePubsubRecvGametvMessage(msg)
		case msg := <-h.subMessages["user.*"]:
			h.handlePubsubRecvUserMessage(msg)
		case msg := <-h.subMessages["conn.*"]:
			h.handlePubsubRecvConnMessage(msg)

		case <-ctx.Done():
			h.log.Debug("hub pubsub ctx done")
			return
		}
	}
}

func (h *Hub) handlePubsubRecvLobbyMessage(m *redis.Message) {
	fmt.Println("hub handlePubsubRecvLobbyMessage", m)

	// ##################################################################
	h.broadcastChannel <- ChannelMessage{channel: "lobby", msg: []byte(m.Payload)}
	// ##################################################################
}

func (h *Hub) handlePubsubRecvGameMessage(m *redis.Message) {
	fmt.Println("hub handlePubsubRecvGameMessage", m)
}

func (h *Hub) handlePubsubRecvGametvMessage(m *redis.Message) {
	fmt.Println("hub handlePubsubRecvGametvMessage", m)
}

func (h *Hub) handlePubsubRecvUserMessage(m *redis.Message) {
	fmt.Println("hub handlePubsubRecvUserMessage", m)

	userID, channel, err := extractUserTopicParts(m.Channel)
	if err != nil {
		return
	}

	h.broadcastUser <- UserMessage{userID: userID, channel: (*Channel)(channel), msg: []byte(m.Payload)}
}

func (h *Hub) handlePubsubRecvConnMessage(m *redis.Message) {
	fmt.Println("hub handlePubsubRecvConnMessage", m)

	connID, err := extractConnTopicParts(m.Channel)
	if err != nil {
		return
	}

	h.broadcastConn <- ConnMessage{connID: connID, msg: []byte(m.Payload)}
}

// extractUserTopicParts extracts the user_id and optional channel if it exists
func extractUserTopicParts(topic string) (uuid.UUID, *string, error) {
	parts := strings.SplitN(topic, ".", 3)
	if len(parts) != 2 && len(parts) != 3 {
		return uuid.Nil, nil, fmt.Errorf("invalid parts length, expected 2 or 3, got: %d", len(parts))
	}

	clientIDStr := parts[1]
	if clientIDStr == "" {
		return uuid.Nil, nil, fmt.Errorf("empty user id")
	}
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		return uuid.Nil, nil, fmt.Errorf("failed to parse user id")
	}
	var channel *string

	if len(parts) == 3 {
		channel = new(parts[2])
	}

	return clientID, channel, nil
}

// extractConnTopicParts extracts the conn_id
func extractConnTopicParts(topic string) (uuid.UUID, error) {
	parts := strings.Split(topic, ".")
	if len(parts) != 2 {
		return uuid.Nil, fmt.Errorf("invalid parts length, expected 2, got: %d", len(parts))
	}
	connIDStr := parts[1]
	if connIDStr == "" {
		return uuid.Nil, fmt.Errorf("empty conn id")
	}
	connID, err := uuid.Parse(connIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse conn id")
	}

	return connID, nil
}

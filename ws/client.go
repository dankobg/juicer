package ws

import (
	"context"
	"log/slog"
	"net/url"
	"sync"
	"time"

	"github.com/coder/websocket"
	pb "github.com/dankobg/juicer/pb/proto/juicer"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
)

type ClientAuthState uint8

const (
	ClientGuest ClientAuthState = iota
	ClientAuth
)

func (cas ClientAuthState) String() string {
	switch cas {
	case ClientAuth:
		return "auth"
	case ClientGuest:
		return "guest"
	default:
		return "unknown"
	}
}

const (
	outMsgBufferSize = 64
	pingPeriod       = 5 * time.Second
	pongWait         = pingPeriod + time.Second*10
)

type client struct {
	userID    uuid.UUID
	connID    uuid.UUID
	conn      *websocket.Conn
	hub       *Hub
	authState ClientAuthState
	channels  []Channel
	outMsg    chan []byte
	query     url.Values
	log       *slog.Logger

	cancel context.CancelFunc

	pongCount    int
	lastPingSent time.Time
	avgLatency   time.Duration
	mu           *sync.RWMutex
}

func NewClient(id uuid.UUID, hub *Hub, conn *websocket.Conn, authState ClientAuthState, logger *slog.Logger, query url.Values,
	cancelFunc context.CancelFunc,
	onPingReceivedRegister func(func(ctx context.Context, payload []byte) bool),
	onPongReceivedRegister func(func(ctx context.Context, payload []byte),
	),
) *client {
	connID := uuid.New()
	clientLogger := logger.With(slog.String("user_id", id.String()), slog.String("conn_id", connID.String()), slog.String("auth_state", authState.String()))

	client := &client{
		userID:    id,
		connID:    connID,
		conn:      conn,
		authState: authState,
		channels:  make([]Channel, 0),
		outMsg:    make(chan []byte, outMsgBufferSize),
		hub:       hub,
		log:       clientLogger,
		query:     query,
		mu:        &sync.RWMutex{},
		cancel:    cancelFunc,
	}

	if onPingReceivedRegister != nil {
		onPingReceivedRegister(client.onPingReceived)
	}

	if onPongReceivedRegister != nil {
		onPongReceivedRegister(client.onPongReceived)
	}

	return client
}

func (c *client) String() string {
	return c.userID.String()
}

func (c *client) JoinChannels(channels []Channel) {
	c.channels = channels
}

func (c *client) ReadLoop(ctx context.Context) {
	defer func() {
		c.hub.ClientDisconnected <- c
	}()

	for {
		msgType, msg, err := c.conn.Read(ctx)
		if err != nil {
			c.log.Debug("conn.Read", slog.Any("error", err))
			return
		}

		c.log.Debug("recv", slog.String("msg_type", msgType.String()), slog.String("msg", string(msg)))

		if err := c.hub.processClientWebsocketMessage(c, msg); err != nil {
			c.log.Error("processClientWebsocketMessage", slog.String("msg_type", msgType.String()), slog.String("msg", string(msg)), slog.Any("error", err))
			// @TODO: some send err msg to client
			continue
		}
	}
}

func (c *client) WriteLoop(ctx context.Context) {
	pingTicker := time.NewTicker(pingPeriod)

	defer func() {
		pingTicker.Stop()
	}()

	for {
		select {
		case outMsg, ok := <-c.outMsg:
			if !ok {
				// c.log.Debug("WriteLoop recv from outMsg chan not ok, closed channel")
				return
			}

			c.log.Debug("conn.Write", slog.String("msg_type", websocket.MessageText.String()), slog.Any("msg", string(outMsg)))

			if err := c.conn.Write(ctx, websocket.MessageText, outMsg); err != nil {
				c.log.Error("conn.Write", slog.Any("error", err))
				return
			}

		case <-pingTicker.C:
			c.mu.Lock()
			c.lastPingSent = time.Now()
			c.mu.Unlock()

			if err := c.conn.Ping(ctx); err != nil {
				c.log.Error("failed to ping", slog.Int("pong_count", c.pongCount), slog.Time("last_ping_sent", c.lastPingSent), slog.Any("error", err))
				return
			}

		case <-ctx.Done():
			c.log.Debug("WriteLoop ctx.Done")
			return
		}
	}
}

func (c *client) onPingReceived(ctx context.Context, payload []byte) bool {
	return true
}

func (c *client) onPongReceived(ctx context.Context, payload []byte) {
	c.mu.RLock()
	curLatency := time.Since(c.lastPingSent)
	c.mu.RUnlock()

	c.mu.Lock()
	c.pongCount++

	// Decaying average after the first four pongs. Stolen from liwords which stole from lichess.
	var mix float64
	if c.pongCount > 4 {
		mix = 0.1
	} else {
		mix = 1 / float64(c.pongCount)
	}

	c.avgLatency += time.Duration(mix * (float64(curLatency) - float64(c.avgLatency)))
	c.mu.Unlock()

	// if c.pongCount%10 == 2 {
	heartbeatMsg := &pb.Message{Event: &pb.Message_Heartbeat{Heartbeat: &pb.Heartbeat{UserId: c.userID.String(), ConnId: c.connID.String(), Guest: c.authState == ClientGuest}}}

	heartbeatMsgBytes, err := protojson.Marshal(heartbeatMsg)
	if err != nil {
		c.log.Error("protojson marshal Message_Heartbeat", slog.Any("error", err))
	} else {
		if err := c.hub.bus.rdb.Publish(context.Background(), "ipc", heartbeatMsgBytes).Err(); err != nil {
			c.log.Error("client publish Message_Heartbeat", slog.String("topic", "ipc"), slog.Any("error", err))
		}
	}
	// }

	latencyMsg := &pb.Message{Event: &pb.Message_Latency{Latency: &pb.Latency{LatencyMs: int32(c.avgLatency.Milliseconds())}}}

	// latencyMsgBytes, err := protojson.Marshal(latencyMsg)
	// if err != nil {
	// 	c.log.Error("protojson marshal Message_Latency", slog.Any("error", err))
	// } else {
	// 	c.outMsg <- latencyMsgBytes
	// }
	_ = latencyMsg
}

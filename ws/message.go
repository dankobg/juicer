package ws

import "github.com/google/uuid"

// ConnMessage is message sent to a single ws conn
type ConnMessage struct {
	connID uuid.UUID
	msg    []byte
}

// UserMessage is message sent to user across all ws conns, unless channel is specified
type UserMessage struct {
	userID  uuid.UUID
	channel *Channel
	msg     []byte
}

// ChannelMessage is message sent to everyone in that channel
type ChannelMessage struct {
	channel Channel
	msg     []byte
}

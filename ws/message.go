package ws

import "github.com/google/uuid"

type ConnMessage struct {
	connID uuid.UUID
	msg    []byte
}

type ClientMessage struct {
	clientID uuid.UUID
	// channel  *Channel
	msg []byte
}

type ChannelMessage struct {
	channel Channel
	msg     []byte
}

package ws

import (
	pb "github.com/dankobg/juicer/pb/proto/juicer"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

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
